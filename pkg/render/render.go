package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("template is not exists")
	}

	/*
		Буферизация `t.Execute(buf, nil)` дает больше контроля над процессом рендеринга чем отправка в поток t.Execute(w, nil)
			- если ошибка происходит в середине рендеринга, пользователь не получит частично сгенерированную страницу (broken html)
			- однократная запись всего содержимого в сетевой поток гораздо эффективней множества мелких
			- постобработка сформированной страницы перед отправкой пользователю
	*/
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, nil)
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser:", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
