package baidu

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type ShareFile struct {
	List    []File `json:"list"`
	ShareId int    `json:"share_id"`
	Uk      int    `json:"uk"`
}

// Verify 获取提取码的sekey 在获取和转存的时候需要用到
func Verify(surl string, pwd string, sekey *string) error {
	param := url.Values{}
	param.Add("surl", surl)
	param.Add("method", "verify")
	data := url.Values{}
	data.Add("pwd", pwd)
	res, err := Api.Client.PostForm(
		Host+"/rest/2.0/xpan/share?"+param.Encode(),
		data,
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("错误的状态码:" + strconv.Itoa(res.StatusCode))
	}
	var body struct {
		resp
		Randsk string `json:"randsk"`
	}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return err
	}
	if body.Errno != 0 {
		return errors.New(body.ErrMsg)
	}
	*sekey = body.Randsk
	return nil
}

// List 获取分享文件列表
func List(surl, sekey string, share *ShareFile) error {
	param := url.Values{}
	param.Add("shorturl", surl)
	param.Add("sekey", sekey)
	param.Add("page", "1")
	param.Add("num", "100")
	param.Add("root", "1")
	param.Add("fid", "0")
	param.Add("method", "list")
	res, err := Api.Client.Get(Host + "/rest/2.0/xpan/share?" + param.Encode())
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("错误的状态码:" + strconv.Itoa(res.StatusCode))
	}
	var body struct {
		resp
		ShareFile
	}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return err
	}
	if body.Errno != 0 {
		return errors.New(body.ErrMsg)
	}
	*share = body.ShareFile
	return nil
}

// Transfer 转存函数
func Transfer(accessToken string, shareId, uk int, sekey, path, fsId string) error {
	param := url.Values{}
	param.Add("method", "transfer")
	param.Add("access_token", accessToken)
	param.Add("shareid", strconv.Itoa(shareId))
	param.Add("from", strconv.Itoa(uk))
	data := url.Values{}
	data.Add("sekey", sekey)
	data.Add("path", path)
	data.Add("fsidlist", fsId)
	res, err := Api.Client.PostForm(
		Host+"/rest/2.0/xpan/share?"+param.Encode(),
		data,
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("错误的状态码:" + strconv.Itoa(res.StatusCode))
	}
	var body struct {
		resp
		FileId string `json:"file_id"`
	}
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return err
	}
	if body.Errno != 0 {
		return errors.New(body.ErrMsg)
	}
	return nil
}
