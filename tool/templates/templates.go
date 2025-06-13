package templates

import _ "embed"

//go:embed page.html
var Page []byte

// TODO(nikonov): consider using two templates

const Intro = `<div class="me">
    Heyo! I'm Alex, a humble software engineer, currently working on my own projects:
    writing an x86 bootloader, learning Linux internals and OS design, and exploring other computer science topics.
    Check out my curiosity playground below.
</div>`
