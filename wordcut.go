package main

import (
	"context"
	"errors"
	"github.com/zone-7/andflow_plugin"
	"regexp"
	"strings"
)
type Andflow_plugin_wordcut struct {

}
func (a *Andflow_plugin_wordcut)GetName() string{
	return "wordsim"
}
func (a *Andflow_plugin_wordcut)Init(callback interface{}){

}
func (a *Andflow_plugin_wordcut)PrepareMetadata(userid int,flowCode string, metadata string)string{
	md:=andflow_plugin.ParseMetadata(metadata)
	md.Name=a.GetName()
	md.Title="文本分词"
	md.Group="自然语言"
	md.Tag="机器学习"
	md.Params=[]andflow_plugin.MetadataPropertiesModel{
		andflow_plugin.MetadataPropertiesModel{Name:"text",Title:"文本参数"},
		andflow_plugin.MetadataPropertiesModel{Name:"dict",Title:"词语列表",Placeholder:"逗号分隔多个词语"},
		andflow_plugin.MetadataPropertiesModel{Name:"result",Title:"分词结果参数"},

	}
	return md.ToJson()

}


func (a *Andflow_plugin_wordcut)Filter(ctx context.Context,runtimeId string,preActionId string, actionId string,callback interface{})(bool,error){

	return true,nil
}


func (a *Andflow_plugin_wordcut)Exec(ctx context.Context,runtimeId string,preActionId string, actionId string,callback interface{})(interface{},error){

	actionCallback := andflow_plugin.ParseActionCallbacker(callback)
	key_text := actionCallback.GetActionParam(actionId, "text")

	dict := actionCallback.GetActionParam(actionId, "dict")

	key_result := actionCallback.GetActionParam(actionId, "result")

	text:=actionCallback.GetRuntimeParam(key_text)
	srcStr,ok :=text .(string)
	if !ok  {
		return nil,errors.New("输入的文本不是字符串")
	}

	g := NewGoJieba()

	if len(dict)>0{
		commonWords:=make([]string,0)
		words := strings.Split(dict,",")
		for _,word:=range words{
			if len(word)==0{
				continue
			}
			w:=strings.Split(word,"，")
			if len(w)==0{
				continue
			}
			commonWords = append(commonWords,w...)
		}

		g.AddWords(commonWords)
	}

 	srcStr = removeHtml(srcStr)
	srcWords := g.C.Cut(srcStr, true)

	actionCallback.SetRuntimeParam(key_result, srcWords)
	return nil,nil
}


func removeHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile(`\\<[\\S\\s]+?\\>`)
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile(`\\<style[\\S\\s]+?\\</style\\>`)
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile(`\\<script[\\S\\s]+?\\</script\\>`)
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile(`\\<[\\S\\s]+?\\>`)
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile(`\\s{2,}`)
	src = re.ReplaceAllString(src, "\n")

	return src
}

func main(){
	andflow_plugin.InitPlugin(&Andflow_plugin_wordcut{})
}