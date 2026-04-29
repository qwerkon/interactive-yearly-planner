{%
{{ if not .Large -}} \renewcommand{\arraystretch}{\myNumArrayStretch}% {{- end}}
\setlength{\tabcolsep}{\myLenTabColSep}%
%
{{ .Month.DefineTable .TableType .Large }}
  {{ if and $.Cfg (not $.Cfg.Pages.MonthlyEnabled) -}}
    {{ .Month.MaybeNameText .Large }}
  {{- else -}}
    {{ .Month.MaybeName .Large }}
  {{- end }}
  {{ if $.Large -}} \hline {{- end }}
  {{ .Month.WeekHeader .Large }} \\ {{ if .Large -}} \noalign{\hrule height \myLenLineThicknessThick} {{- else -}} \hline {{- end}}
  {{- range $i, $week := .Month.Weeks }}
  {{if $.Cfg}}{{if $.Cfg.Pages.WeeklyEnabled}}{{$week.WeekNumber $.Large}}{{else if $.Cfg.Pages.DailyEnabled}}{{$week.WeekNumberDay $.Large}}{{else}}{{$week.WeekNumberText $.Large}}{{end}}{{else}}{{$week.WeekNumber $.Large}}{{end}} &
    {{- range $j, $day := $week.Days -}}
      {{- if and $.Cfg (not $.Cfg.Pages.DailyEnabled) $.Cfg.Pages.WeeklyEnabled -}}
        {{- $day.Day $.Today $.Large $.Cfg ($week.Ref) -}}
      {{- else -}}
        {{- $day.Day $.Today $.Large $.Cfg -}}
      {{- end -}}
      {{- if eq $j 6 -}}
        \\ {{ if $.Large -}} \hline {{- end -}}
      {{- else -}} & {{- end -}}
    {{- end -}}
  {{ end }}
  {{ .Month.EndTable .TableType -}}
}
