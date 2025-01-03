package tools

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utilHttp"
	"github.com/hilaoyu/password-manager/config"
	"github.com/hilaoyu/password-manager/ui"
	"image/color"
	"strconv"
	"strings"
	"time"
)

func ToolAesRequest() (content *fyne.Container) {

	methodInput := widget.NewEntry()
	urlInput := widget.NewEntry()
	secretInput := widget.NewEntry()
	appIdInput := widget.NewEntry()

	paramsInput := widget.NewMultiLineEntry()
	paramsInput.Wrapping = fyne.TextWrapBreak

	dataInput := widget.NewMultiLineEntry()
	dataInput.Wrapping = fyne.TextWrapBreak

	headersInput := widget.NewMultiLineEntry()
	headersInput.Wrapping = fyne.TextWrapBreak
	headersInput.Text = "{\"Content-Type\":\"application/json\"}"

	requestDataInput := widget.NewMultiLineEntry()
	requestDataInput.Wrapping = fyne.TextWrapBreak

	resultBodyInput := widget.NewMultiLineEntry()
	resultBodyInput.Wrapping = fyne.TextWrapBreak

	resultDataDecodeInput := widget.NewMultiLineEntry()
	resultDataDecodeInput.Wrapping = fyne.TextWrapBreak

	requestButton := widget.NewButton("发送请求", func() {
		methodInput.Text = strings.TrimSpace(methodInput.Text)
		if "" == methodInput.Text {
			methodInput.Text = "get"
			methodInput.Refresh()
		}
		urlInput.Text = strings.TrimSpace(urlInput.Text)
		if "" == urlInput.Text {
			config.UiDefault().WindowError(fmt.Errorf("地址不能为空"))
			return
		}
		secretInput.Text = strings.TrimSpace(secretInput.Text)
		if "" == secretInput.Text {
			config.UiDefault().WindowError(fmt.Errorf("密钥不能为空"))
			return
		}

		var err error
		params := map[string]interface{}{}
		paramsInput.Text = strings.TrimSpace(paramsInput.Text)
		if "" != paramsInput.Text {
			err = json.Unmarshal([]byte(paramsInput.Text), &params)
			if nil != err {
				config.UiDefault().WindowError(fmt.Errorf("参数错误: %v", err))
				return
			}
		}

		dataInput.Text = strings.TrimSpace(dataInput.Text)
		if "" == dataInput.Text {
			config.UiDefault().WindowError(fmt.Errorf("数据不能为空"))
			return
		}

		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(dataInput.Text), &data)
		if nil != err {
			config.UiDefault().WindowError(fmt.Errorf("数据错误: %v", err))
			return
		}

		headers := map[string]string{}
		headersInput.Text = strings.TrimSpace(headersInput.Text)
		if "" != headersInput.Text {
			err = json.Unmarshal([]byte(headersInput.Text), &headers)
			if nil != err {
				config.UiDefault().WindowError(fmt.Errorf("请求头错误: %v", err))
				return
			}
		}

		encryptor := utilEnc.NewAesEncryptor(secretInput.Text)
		data["_timestamp"] = time.Now().UTC().Unix()
		data["_data_id"] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		enData, err := encryptor.Encrypt(data)
		params["data"] = enData

		appIdInput.Text = strings.TrimSpace(appIdInput.Text)
		if "" != appIdInput.Text {
			params["app_id"] = appIdInput.Text
		}

		paramsJsonByte, err := json.Marshal(params)
		if nil != err {
			config.UiDefault().WindowError(fmt.Errorf("请求数据转json错误: %v", err))
			return
		}
		paramsStr := string(paramsJsonByte)

		requestDataInput.Text = paramsStr
		requestDataInput.Refresh()

		httpClient := utilHttp.NewHttpClient(urlInput.Text)

		httpClient.WithRawBody(paramsStr)

		body, err := httpClient.RequestPlain(methodInput.Text, "", headers)
		if nil != err {
			config.UiDefault().WindowError(fmt.Errorf("请求出错: %v", err))
			return
		}

		resultBodyInput.Text = string(body)
		resultBodyInput.Refresh()

		result := map[string]interface{}{}

		err = json.Unmarshal(body, &result)
		if nil != err {
			config.UiDefault().WindowError(fmt.Errorf("解析body: %v", err))
			return
		}

		if resultEnData, ok := result["data"]; ok {
			if resultEnDataStr, ok1 := resultEnData.(string); ok1 {
				resultDeDataStr, err := encryptor.DecryptString(resultEnDataStr)
				if nil != err {
					config.UiDefault().WindowError(fmt.Errorf("解密结果错误: %v", err))
					return
				}

				resultDataDecodeInput.Text = resultDeDataStr
				resultDataDecodeInput.Refresh()
			}
		}

	})

	form := widget.NewForm(
		widget.NewFormItem("请求方式:", methodInput),
		widget.NewFormItem("请求地址:", urlInput),
		widget.NewFormItem("密钥:", secretInput),
		widget.NewFormItem("应用ID:", appIdInput),

		widget.NewFormItem("加密数据:", dataInput),
		widget.NewFormItem("明文数据:", paramsInput),
		widget.NewFormItem("请求头:", headersInput),

		widget.NewFormItem("", requestButton),
		widget.NewFormItem("", ui.NewRectangleWithSize(color.Transparent, 5, 0)),

		widget.NewFormItem("请求内容:", requestDataInput),
		widget.NewFormItem("返回内容:", resultBodyInput),
		widget.NewFormItem("返回解密:", resultDataDecodeInput),
	)

	content = container.NewStack(ui.NewScrollWithSize(form, 600, 800))
	return
}
