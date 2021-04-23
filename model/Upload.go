package model

import (
	"context"
	"ginVue/utils/errmsg"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
)

// 控制数据的写入和操作一般写在service层下面

// 包级变量
var AccessKey = "tnF4cmhr1A4VeiFvoM9eu2wg8Xw7BOECifuQUeTz"
var SecretKey = "3Drh5iB87fUH7hkRnpR1mweq-DglcSKOdPCG4fGI"
var Bucket = "bai1889"
var ImgUrl = "http://qr8a78yya.hn-bkt.clouddn.com/"


func UpLoadFile(file multipart.File, fileSize int64) (string, int) {
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", errmsg.ERROR
	}
	url := ImgUrl + ret.Key
	return url, errmsg.SUCCESS

}