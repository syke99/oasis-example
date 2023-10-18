package main

var hello = `<div id={{ .props.ID }} class={{ .props.classes }} hx-get={{ .props.url }} hx-target={{ .props.target }} hx-swap={{ .props.swap }} hx-trigger={{ .props.trigger }}></div>`

var hi = `<p>Hello {{ .props.person.Name }}!!{{if .payload.isBirthday }} Today is your birthday; you're now {{ .payload.age }} years old!!{{end}}</p>`
