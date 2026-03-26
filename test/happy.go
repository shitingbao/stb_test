package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
)

type HappyData struct {
	Red string `json:"red"`
}

type HappyResult struct {
	Tflag    int         `json:"Tflag"`
	Message  string      `json:"message"`
	PageNo   int         `json:"pageNo"`
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
	State    int         `json:"state"`
	Total    int         `json:"total"`
	Result   []HappyData `json:"result"`
}

// 快乐8
// 接口数据
// curl --location --request GET 'https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=kl8&issueCount=&issueStart=&issueEnd=&dayStart=&dayEnd=&pageNo=2&pageSize=30&week=&systemType=PC'\” \   -H '\”Accept: application/json, text/javascript, */*; q=0.01'\” \   -H '\”Accept-Language: zh-CN,zh;q=0.9'\” \   -H '\”Cache-Control: no-cache'\” \   -H '\”Connection: keep-alive'\” \   -b '\”HMF_CI=73d8c9f0fb03dc6f6122b48b3a127ecc01e54b5e5cd61ec789780fa3bed9e662820bf52721c92e55feab40a130bffcfea7fe5779bd813823b0eedf16c12fa639b5; 21_vq=1'\” \   -H '\”Pragma: no-cache'\” \   -H '\”Referer: https://www.cwl.gov.cn/ygkj/wqkjgg/kl8/'\” \   -H '\”Sec-Fetch-Dest: empty'\” \   -H '\”Sec-Fetch-Mode: cors'\” \   -H '\”Sec-Fetch-Site: same-origin'\” \   -H '\”User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36'\” \   -H '\”X-Requested-With: XMLHttpRequest'\” \   -H '\”sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"'\” \   -H '\”sec-ch-ua-mobile: ?0'\” \   -H '\”sec-ch-ua-platform: "macOS"'\”&systemType=PC' \
// --header 'Accept: application/json, text/javascript, */*; q=0.01' \
// --header 'Host: www.cwl.gov.cn' \
// --header 'Referer: https://www.cwl.gov.cn/ygkj/wqkjgg/kl8/' \
// --header 'X-Requested-With: XMLHttpRequest' \
// --header 'Cookie: HMF_CI=73d8c9f0fb03dc6f6122b48b3a127ecc01e54b5e5cd61ec789780fa3bed9e662820bf52721c92e55feab40a130bffcfea7fe5779bd813823b0eedf16c12fa639b5; 21_vq=1' \
// --header 'User-Agent: Apifox/1.0.0 (https://apifox.com)'
func getHappy8(pageNo, pageSize int) (HappyResult, error) {
	const endpoint = "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice"

	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 30
	}

	q := url.Values{}
	q.Set("name", "kl8")
	q.Set("issueCount", "")
	q.Set("issueStart", "")
	q.Set("issueEnd", "")
	q.Set("dayStart", "")
	q.Set("dayEnd", "")
	q.Set("pageNo", strconv.Itoa(pageNo))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("week", "")
	q.Set("systemType", "PC")

	// 该站点对 Go 的默认 TLS/UA 指纹可能直接返回 403，因此这里用系统 curl 发起请求更稳定。
	// 如需最新 Cookie，可设置环境变量 HAPPY_COOKIE 覆盖。
	ck := os.Getenv("HAPPY_COOKIE")
	if ck == "" {
		ck = "HMF_CI=73d8c9f0fb03dc6f6122b48b3a127ecc01e54b5e5cd61ec789780fa3bed9e662820bf52721c92e55feab40a130bffcfea7fe5779bd813823b0eedf16c12fa639b5; 21_vq=1"
	}

	cmd := exec.Command(
		"curl",
		"-sSL",
		endpoint+"?"+q.Encode(),
		"-H", "Accept: application/json, text/javascript, */*; q=0.01",
		"-H", "Accept-Language: zh-CN,zh;q=0.9",
		"-H", "Cache-Control: no-cache",
		"-H", "Pragma: no-cache",
		"-H", "Referer: https://www.cwl.gov.cn/ygkj/wqkjgg/kl8/",
		"-H", "X-Requested-With: XMLHttpRequest",
		"-H", "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
		"-H", "sec-ch-ua: \"Chromium\";v=\"146\", \"Not-A.Brand\";v=\"24\", \"Google Chrome\";v=\"146\"",
		"-H", "sec-ch-ua-mobile: ?0",
		"-H", "sec-ch-ua-platform: \"macOS\"",
		"-H", "Sec-Fetch-Dest: empty",
		"-H", "Sec-Fetch-Mode: cors",
		"-H", "Sec-Fetch-Site: same-origin",
		"-H", "Cookie: "+ck,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return HappyResult{}, fmt.Errorf("curl failed: %w; stderr=%s", err, stderr.String())
	}
	body := stdout.Bytes()

	var r HappyResult
	if err := json.Unmarshal(body, &r); err != nil {
		return HappyResult{}, fmt.Errorf("unmarshal failed: %w; body=%s", err, string(body))
	}
	return r, nil
}

func loadHappy(startPage int) (int, error) {
	if startPage <= 0 {
		startPage = 1
	}

	pageSize := 30
	pageNo := startPage
	totalFetched := 0

	f, err := os.OpenFile("data.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()

	for {
		r, err := getHappy8(pageNo, pageSize)
		if err != nil {
			return totalFetched, err
		}

		if len(r.Result) == 0 {
			break
		}

		for _, item := range r.Result {
			if _, err := f.WriteString(item.Red + "\n"); err != nil {
				return totalFetched, err
			}
		}

		totalFetched += len(r.Result)
		fmt.Printf("pageNo=%d fetched=%d (total=%d pageNum=%d)\n", r.PageNo, len(r.Result), r.Total, r.PageNum)

		// 1) 优先用 pageNum（如果它代表总页数）
		if r.PageNum > 0 && r.PageNo >= r.PageNum {
			break
		}
		// 2) 再用 total 推断是否到最后一页
		if r.Total > 0 && r.PageNo*pageSize >= r.Total {
			break
		}

		pageNo++
	}

	return totalFetched, nil
}
