package templates

import _ "embed"

//go:embed page.html
var Page []byte

const Intro = `<p class="me">
Heyo! I'm Alex, a humble software engineer, currently working at Tabby.
I'm interested in how systems work under the hood, so I spend my time exploring Linux, networking, and other low-level stuff.
Check out my curiosity playground below.</p>`
