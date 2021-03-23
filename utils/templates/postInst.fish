go mod download -x
fisher install jethrokuan/fzf
fisher install jhillyerd/plugin-git
{{ range . }}
fisher install {{ . }}
{{ end }}