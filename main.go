package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"
)

func main() {
	// Backend Server
	bs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("---------------- Backend Server ------------------")
		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			fmt.Printf("Backend httputil.DumpRequest failed: %w\n", err)
			return
		}
		fmt.Printf(string(dump))
		fmt.Println("--- Backend Server が受信した RemoteAddr ---")
		fmt.Printf("RemoteAddr: %s\n", r.RemoteAddr)
		fmt.Println()

		ip := r.RemoteAddr
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			ips := strings.Split(xff, ",")
			ip = strings.TrimSpace(ips[0])
		}

		fmt.Println("--- 最終的な X-Forwarded-For と送信元IP ---")
		fmt.Printf("X-Forwarded-For: %s\n", xff)
		fmt.Printf("Source IP: %s\n", ip)
		fmt.Println()

		if ip == "127.0.0.1" {
			f, err := os.Open("flag.txt")
			if err != nil {
				fmt.Printf("os.Open failed: %v\n", err)
				return
			}
			defer f.Close()

			answer, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Printf("ioutil.ReadAll failed: %v\n", err)
				return
			}
			fmt.Printf("正解！秘密の答えをレスポンスしました！\n")
			fmt.Fprintf(w, string(answer))
		} else {
			fmt.Printf("残念！送信元IP が 127.0.0.1 になるようにリクエストを送ってください！(%s != 127.0.0.1)\n", ip)
			fmt.Fprintf(w, "もっと頑張りましょう！")
		}

		fmt.Println("==========================")
	}))
	defer bs.Close()

	// Middle Proxy Server
	mp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("---------------- Middle Proxy ----------------")

		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			fmt.Printf("Middle Proxy httputil.DumpRequest failed: %v\n", err)
			return
		}
		fmt.Printf(string(dump))
		fmt.Println("--- Middle Proxy が受信した RemoteAddr ---")
		fmt.Printf("RemoteAddr: %s\n", r.RemoteAddr)
		fmt.Println()

		burl, err := url.Parse(bs.URL)
		if err != nil {
			fmt.Printf("Front Proxy url.Parse failed: %v\n", err)
			return
		}
		httputil.NewSingleHostReverseProxy(burl).ServeHTTP(w, r)
	}))
	defer mp.Close()

	// Front Proxy Server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("---------------- Front Proxy ----------------")
		r.Header.Del("X-Forwarded-For")

		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			fmt.Printf("Front Proxy httputil.DumpRequest failed: %v\n", err)
			return
		}
		fmt.Printf(string(dump))
		fmt.Println("--- Front Proxy が受信した RemoteAddr ---")
		fmt.Printf("RemoteAddr: %s\n", r.RemoteAddr)
		fmt.Println()

		murl, err := url.Parse(mp.URL)
		if err != nil {
			fmt.Printf("Front Proxy url.Parse failed: %v\n", err)
			return
		}
		httputil.NewSingleHostReverseProxy(murl).ServeHTTP(w, r)
	})

	// Run server
	fmt.Printf("========================================================\n")
	fmt.Printf("Welcome to Kanmu Office hour @ Go Conference 2021 Autumn\n")
	fmt.Printf(`"Go" beyond your proxy! `)
	fmt.Printf("Go version: %s\n", runtime.Version())
	http.ListenAndServe(":8000", nil)
}
