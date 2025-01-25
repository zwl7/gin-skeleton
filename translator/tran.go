package translator

import (
	"fmt"
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

func Validate(data interface{}) (string,int) {
	//验证对象 并将英文报错转换成中文报错（message）
	validate := validator.New()
	uni := unTrans.New(zh_Hans_CN.New())
	trans, _ := uni.GetTranslator("zh_Hans_CN")
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println("err:", err)
	}
	//将验证法字段名 映射为中文名
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})
	err = validate.Struct(data)
	if err != nil {
		//错误可能有多个 遍历 返回一个
		for _, v := range err.(validator.ValidationErrors) {
			return v.Translate(trans), 500
		}
	}
	return "", 200
}
