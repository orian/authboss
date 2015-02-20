package render

import (
	"bytes"
	"io"
	"net/http"

	"gopkg.in/authboss.v0"
	"gopkg.in/authboss.v0/internal/views"
)

// View renders a view with xsrf and flash attributes.
func View(ctx *authboss.Context, w http.ResponseWriter, r *http.Request, t views.Templates, name string, data authboss.HTMLData) error {
	tpl, ok := t[name]
	if !ok {
		return authboss.RenderErr{tpl.Name(), data, ErrTemplateNotFound}
	}

	data.Merge("xsrfName", authboss.Cfg.XSRFName, "xsrfToken", authboss.Cfg.XSRFMaker(w, r))

	if flash, ok := ctx.CookieStorer.Get(authboss.FlashSuccessKey); ok {
		ctx.CookieStorer.Del(authboss.FlashSuccessKey)
		data.Merge(authboss.FlashSuccessKey, flash)
	}
	if flash, ok := ctx.CookieStorer.Get(authboss.FlashErrorKey); ok {
		ctx.CookieStorer.Del(authboss.FlashErrorKey)
		data.Merge(authboss.FlashErrorKey, flash)
	}

	if authboss.Cfg.LayoutDataMaker != nil {
		data.Merge(authboss.Cfg.LayoutDataMaker(w, r))
	}

	buffer = &bytes.Buffer{}
	err = tpl.ExecuteTemplate(buffer, tpl.Name(), data)
	if err != nil {
		return authboss.RenderErr{tpl.Name(), data, err}
	}

	err = io.Copy(w, buffer)
	if err != nil {
		return authboss.RenderErr{tpl.Name(), data, err}
	}
}

// Redirect sets any flash messages given and redirects the user.
func Redirect(ctx *authboss.Context, w http.ResponseWriter, r *http.Request, path, flashSuccess, flashError string) {
	if len(flashSuccess) > 0 {
		ctx.CookieStorer.Put(authboss.FlashSuccessKey, flashSuccess)
	}
	if len(flashError) > 0 {
		ctx.CookieStorer.Put(authboss.FlashErrorKey, flashError)
	}
	http.Redirect(w, r, path, http.StatusTemporaryRedirect)
}