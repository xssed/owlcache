package network

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
)

// 发起请求获取URL数据并进行缓存处理
func (owlhandler *OwlHandler) GetUrlCacheData(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, []byte) {

	//先去查询内存数据库中是否存在这个值
	//设置Key值,定位到Uri
	key := owlhandler.owlrequest.Key
	//判断是否开启Urlcache的快捷访问
	if owlconfig.OwlConfigModel.Open_Urlcache == "1" && owlconfig.OwlConfigModel.Urlcache_Request_Easy == "1" && len(owlhandler.owlrequest.Key) < 1 {
		//开启Urlcache的快捷访问后重新定义key值
		key = r.RequestURI
		owlhandler.owlrequest.Key = key
	}
	//执行K/V数据查询,本地内存数据库->Memcache(如果开启)->Redis(如果开启）
	owlhandler.baseget()
	//定义返回信息变量
	var print []byte
	//执行K/V数据查询后，查询到具有缓存数据
	if owlhandler.owlresponse.Data != nil {
		w.WriteHeader(200)
		print = owlhandler.owlresponse.Data
		return w, print
	} else {
		//未查询到数据后去请求URL数据缓存到本地
		return owlhandler.GetUrlCacheExe(w, r, print)
	}

}

// 未查询到数据后去请求URL数据缓存到本地
func (owlhandler *OwlHandler) GetUrlCacheExe(w http.ResponseWriter, r *http.Request, print []byte) (http.ResponseWriter, []byte) {

	//站点索引值
	var site_index *int
	//处理查找的站点索引值
	//未设置初始值，错误参数，超出索引边界都将设置为默认值0
	if (r.FormValue("uc_site")) != "" {
		temp_index, cerr := strconv.Atoi(r.FormValue("uc_site"))
		//参数错误,超出索引边界
		if cerr != nil || temp_index+1 > len(owlconfig.OwlUCConfigModel.SiteList) {
			temp_index = 0
		}
		site_index = &temp_index
	} else {
		var temp_index int = 0
		site_index = &temp_index
	}

	//使用gorequest类库发起http client请求获取数据
	resp, errs := owlhandler.getUrlDataExe(site_index, owlhandler.owlrequest.Key, r)
	//有错误,报出异常
	if errs != nil {
		w.WriteHeader(500)
		print = []byte(errs.Error())
		return w, print
	}

	//数据处理成功，内容赋值部分
	owlhandler.Transmit(SUCCESS)
	owlhandler.owlresponse.Key = owlhandler.owlrequest.Key
	owlhandler.owlresponse.Data = resp.Byte_slices
	print = owlhandler.owlresponse.Data
	w.WriteHeader(200)
	return w, print

}

func (owlhandler *OwlHandler) getUrlDataExe(index *int, key string, r *http.Request) (UCtext, error) {

	//使用SingleFlight来限制高并发场景下对同一个后端接口的并发量，减缓后端接口压力
	value, _, _ := SingleFlightGroupUC.Do(key, func() (ret interface{}, err error) {

		uct := owlhandler.getUrlData(index, key, r)

		defer uct.Res.Body.Close() //资源释放

		//判断http client请求获取数据是否请求异常
		if uct.Err_slices != nil || uct.Res.StatusCode >= 300 {
			errstr := owltools.ErrorSliceJoinToString(uct.Err_slices)
			if errstr != "" {
				errlog := owltools.JoinString("getUrlDataExe method getUrlData error:", errstr)
				owllog.OwlLogUC.Info(errlog) //日志记录
				return uct, errors.New(errlog)
			}
		}
		//读取获取到的数据
		body, ioerr := io.ReadAll(uct.Res.Body)
		//判断响应数据读取异常
		if ioerr != nil {
			ioutil_errlog := owltools.JoinString("getUrlDataExe method getUrlData io.ReadAll error:", ioerr.Error())
			owllog.OwlLogUC.Info(ioutil_errlog) //日志记录
			return uct, errors.New(ioutil_errlog)
		}
		//判断是否获取到有效数据
		if uint64(len(body)) > 1 {
			//将数据存储到内存数据库
			//设置站点配置信息
			site := owlconfig.OwlUCConfigModel.SiteList[*index]
			//先判断是否要验证tonken
			if site.CheckToken == 1 {
				//验证未通过，不存储数据
				if !owlhandler.CheckAuth(r) {
					un_auth_log := owltools.JoinString("Key:", owlhandler.owlrequest.Key, " Token verification is not pass, no write in database.")
					owllog.OwlLogUC.Info(un_auth_log) //日志记录
					return uct, errors.New(un_auth_log)
				}
			}
			//存储信息
			owlhandler.owlrequest.Key = key
			owlhandler.owlrequest.Value = body
			exptime := time.Duration(site.KeyExpire) * time.Second
			owlhandler.owlrequest.Expires = exptime
			owlhandler.Set()
			//返回数据信息
			return uct, nil
		}
		return uct, errors.New("getUrlDataExe method getUrlData Data is empty!")
	})

	return_uct, _ := value.(UCtext)
	return return_uct, nil

}

// 使用gorequest类库发起http client请求获取数据
// 请求 站点的索引值:index,需要查找的key值或者Uri:key,*http.Request
// 返回 gorequest.Response，http响应的byte数据，http请求的错误信息
func (owlhandler *OwlHandler) getUrlData(index *int, key string, r *http.Request) UCtext {

	//设置站点配置
	site := owlconfig.OwlUCConfigModel.SiteList[*index]
	//定义url地址
	var url string = owltools.JoinString(site.Host, key)
	//创建http client
	var grsa *gorequest.SuperAgent
	grsa = gorequest.New()
	//设置请求地址
	var r_method string
	if r.Method == "POST" {
		r_method = "POST"
		grsa.Post(url)
		//将POST请求数据转发给目标
		r.ParseForm() //解析参数
		for p_k, p_v := range r.Form {
			grsa.Param(p_k, p_v[0])
		}
	} else {
		r_method = "GET"
		grsa.Get(url)
	}
	//查询是否设置代理服务器
	if len(site.Proxy) > 7 {
		grsa.Proxy(site.Proxy)
	}
	//设置超时
	grsa.Timeout(time.Duration(site.Timeout) * time.Millisecond)
	//Headers,如果Request headers的子项有设置值
	if len(site.Headers) != 0 {
		//遍历自定义超集集合
		for _, rhs := range site.Headers {
			//遍历取出要发送的值
			for _, rh := range rhs.Value {
				//首字母进行大写
				rh = owltools.Ucfirst(rh)
				rhv, exist := r.Header[rh] //从Request headers中取值
				if exist {
					grsa.AppendHeader(rh, owltools.StringSliceJoinToString(rhv))
				}
			}
		}
	}

	owllog.OwlLogUC.Info("UCClient URL:", url, " Method:", r_method, " Request Headers:", grsa.Header, " QueryData:", grsa.QueryData) //请求的日志记录

	//发送请求获取数据
	r_res, r_byte_slices, r_err_slices := grsa.EndBytes()
	//判断请求的响应数据是否超过本地允许的最大值
	if r_err_slices == nil {
		if uint64(len(r_byte_slices)) > site.MaxStorageLimit {
			r_err_slices = append(r_err_slices, ErrorUCMaxStorageLimitOver)
		}
	}

	//清理资源
	grsa.ClearSuperAgent()

	return UCtext{Res: r_res, Byte_slices: r_byte_slices, Err_slices: r_err_slices}

}

type UCtext struct {
	Res         gorequest.Response
	Byte_slices []byte
	Err_slices  []error
}
