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
  {{if and $.Cfg (not $.Cfg.Pages.WeeklyEnabled)}}{{$week.WeekNumberText $.Large}}{{else}}{{$week.WeekNumber $.Large}}{{end}} &
    {{- range $j, $day := $week.Days -}}
      {{- $day.Day $.Today $.Large -}}
      {{- if eq $j 6 -}}
        \\ {{ if $.Large -}} \hline {{- end -}}
      {{- else -}} & {{- end -}}
    {{- end -}}
  {{ end }}
  {{ .Month.EndTable .TableType -}}
}
