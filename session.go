package toold

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Session 对象
type Session interface {
	Set(key, value interface{}) error //设置Session
	Get(key interface{}) interface{}  //获取Session
	Delete(key interface{}) error     //删除Session
	SessionID() string                //当前SessionID
}

//Provider 管理器
type Provider interface {
	SessionInit(sid string) (Session, error) //初始化
	SessionRead(sid string) (Session, error) //返回由相应 sid 表示的 Session
	SessionDestroy(sid string) error         //给定一个 sid，删除相应的 Session
	SessionGC(maxLifeTime int64)             // 删除过期的 Session 变量
}

var providers = make(map[string]Provider)

//RegisterProvider 注册一个能通过名称来获取的 session provider 管理器
func RegisterProvider(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}

	if _, p := providers[name]; p {
		panic("session: Register provider is existed")
	}

	providers[name] = provider
}

//Manager 全局的 Session 的管理器
type Manager struct {
	cookieName  string     //cookie的名称
	lock        sync.Mutex //锁，保证并发时数据的安全一致
	provider    Provider   //管理session
	maxLifeTime int64      //超时时间
}

//NewManager 创建
func NewManager(providerName, cookieName string, maxLifetime int64) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", providerName)
	}

	//返回一个 Manager 对象
	return &Manager{
		cookieName:  cookieName,
		maxLifeTime: maxLifetime,
		provider:    provider,
	}, nil
}

var globalSession *Manager

func init() {
	globalSession, _ = NewManager("memory", "sessionid", 3600)
}

func (manager *Manager) sessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

//SessionStart 根据当前请求的cookie中判断是否存在有效的session, 不存在则创建
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {

	fmt.Printf("sidsd:%v", manager)
	//为该方法加锁
	manager.lock.Lock()
	defer manager.lock.Unlock()
	//获取 request 请求中的 cookie 值
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionID()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(sid), //转义特殊符号@#￥%+*-等
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie) //将新的cookie设置到响应中
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

// SessionDestroy 注销 Session
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}

	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.provider.SessionDestroy(cookie.Value)
	expiredTime := time.Now()
	newCookie := http.Cookie{
		Name: manager.cookieName,
		Path: "/", HttpOnly: true,
		Expires: expiredTime,
		MaxAge:  -1, //会话级cookie
	}
	http.SetCookie(w, &newCookie)
}

//SessionGC 删除session
func (manager *Manager) SessionGC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	//使用time包中的计时器功能，它会在session超时时自动调用GC方法
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() {
		manager.SessionGC()
	})
}

//记录该session被访问的次数
func count(w http.ResponseWriter, r *http.Request) {
	sess := globalSession.SessionStart(w, r) //获取session实例
	createTime := sess.Get("createTime")     //获得该session的创建时间
	if createTime == nil {
		sess.Set("createTime", time.Now().Unix())
	} else if (createTime.(int64) + 360) < (time.Now().Unix()) { //已过期
		//注销旧的session信息，并新建一个session  globalSession.SessionDestroy(w, r)
		sess = globalSession.SessionStart(w, r)
	}
	count := sess.Get("countnum")
	if count == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", count.(int)+1)
	}
}
