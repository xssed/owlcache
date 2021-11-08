package network

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
)

//发起请求获取URL数据并进行缓存处理
func (owlhandler *OwlHandler) GeUrlCacheData(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, []byte) {

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

		//站点索引值
		var site_index int = 0
		//处理查找的站点索引值
		//未设置初始值，错误参数，超出索引边界都将设置为默认值0
		if (r.FormValue("uc_site")) != "" {
			site_index, cerr := strconv.Atoi(r.FormValue("uc_site"))
			//参数错误,超出索引边界
			if cerr != nil || site_index+1 > len(owlconfig.OwlUCConfigModel.SiteList) {
				site_index = 0
			}
		}
		//使用gorequest类库发起http client请求获取数据
		resp, _, errs := owlhandler.getUrlData(site_index, key, r)
		if errs != nil || resp.StatusCode >= 300 {
			errstr := owltools.ErrorSliceJoinToString(errs)
			if errstr != "" {
				owllog.OwlLogUC.Info(owltools.JoinString("GeUrlCacheData method getUrlData error:", errstr)) //日志记录
				//http client请求获取数据，请求异常
				w.WriteHeader(500)
				print = []byte(errstr)
				return w, print
			}
		}

		defer resp.Body.Close() //资源释放
		body, ioerr := ioutil.ReadAll(resp.Body)
		if ioerr != nil {
			owllog.OwlLogUC.Info(owltools.JoinString("GeUrlCacheData method getUrlData ioutil.ReadAll error:", ioerr.Error())) //日志记录
			//响应数据读取异常
			w.WriteHeader(500)
			print = []byte(ioerr.Error())
			return w, print
		}
		// out, _ := os.Create("io.jpg")
		// io.Copy(out, bytes.NewReader(body))

		//数据处理成功，内容赋值部分
		owlhandler.Transmit(SUCCESS)
		owlhandler.owlresponse.Key = key
		owlhandler.owlresponse.Data = body

		if owlhandler.owlresponse.Data != nil {
			w.WriteHeader(200)
			print = owlhandler.owlresponse.Data
			//将数据存储到内存数据库
			//设置站点配置信息
			site := owlconfig.OwlUCConfigModel.SiteList[site_index]
			//先判断是否要验证tonken
			if site.CheckToken == 1 {
				if !owlhandler.CheckAuth(r) {
					owllog.OwlLogUC.Info(owltools.JoinString("Key:", key, " Token verification is not pass, no write in database.")) //日志记录
					//验证未通过，不存储数据，直接返回数据信息
					return w, print
				}
			}
			//存储信息
			owlhandler.owlrequest.Value = print
			exptime := time.Duration(site.KeyExpire) * time.Second
			owlhandler.owlrequest.Expires = exptime
			owlhandler.Set()
			//返回数据信息
			return w, print
		} else {
			print = []byte("GeUrlCacheData method getUrlData Data is empty!")
			return w, print
		}

	}

}

//使用gorequest类库发起http client请求获取数据
//请求 站点的索引值:index,需要查找的key值或者Uri:key,*http.Request
//返回 gorequest.Response，http响应的byte数据，http请求的错误信息
func (owlhandler *OwlHandler) getUrlData(index int, key string, r *http.Request) (gorequest.Response, []byte, []error) {

	//设置站点配置
	site := owlconfig.OwlUCConfigModel.SiteList[index]
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
		// ct := r_res.Header.Get("Content-Length")
		// if ct != "" {
		// 	rct, parerr := strconv.ParseUint(ct, 0, 64)
		// 	if parerr != nil {
		// 		fmt.Println(parerr)
		// 	}
		// 	if rct > site.MaxStorageLimit {
		// 		r_err_slices = append(r_err_slices, ErrorUCMaxStorageLimitOver)
		// 	}
		// }
		if uint64(len(r_byte_slices)) > site.MaxStorageLimit {
			r_err_slices = append(r_err_slices, ErrorUCMaxStorageLimitOver)
		}
	}

	//清理资源
	grsa.ClearSuperAgent()

	return r_res, r_byte_slices, r_err_slices

}
