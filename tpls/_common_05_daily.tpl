{{- $today := .Body.Day -}}

\begin{minipage}[t]{ {{- if .Cfg.Layout.Lengths.DailyRightWidth -}} \dimexpr\linewidth-{{ .Cfg.Layout.Lengths.DailyRightWidth }}-\myLenTriColSep {{- else -}} \myLenTriCol {{- end -}} }
{{template "schedule.tpl" dict "Cfg" .Cfg "Day" .Body.Day}}
  \vspace{\dimexpr4mm+.3pt}

{{- if .Cfg.CalAfterSchedule -}}
{{- template "monthTabularV2.tpl" dict "Cfg" .Cfg "Month" .Body.Month "Today" $today -}}
{{- end -}}
\end{minipage}%
\hspace{\myLenTriColSep}%
\begin{minipage}[t]{ {{- if .Cfg.Layout.Lengths.DailyRightWidth -}} {{ .Cfg.Layout.Lengths.DailyRightWidth }} {{- else -}} \dimexpr2\myLenTriCol+\myLenTriColSep {{- end -}} }
  \myUnderline{Top priorities\myDummyQ}
  \Repeat{\myNumDailyTodos}{\myTodoLineGray}
  \vskip\dimexpr5.4mm
  \myUnderline{Notes{{ if or .Cfg.Pages.DailyNotesEnabled .Cfg.Pages.DailyReflectEnabled .Cfg.Layout.Numbers.NotesIndexPages }} $\vert$ {{ end }}{{ if .Cfg.Pages.DailyNotesEnabled }}{{ $today.LinkLeaf "More" "More" }}{{ end }}{{ if and .Cfg.Pages.DailyNotesEnabled .Cfg.Pages.DailyReflectEnabled }}\hfill{}{{ end }}{{ if .Cfg.Pages.DailyReflectEnabled }}{{ $today.LinkLeaf "Reflect" "Reflect" }}{{ end }}{{ if and .Cfg.Layout.Numbers.NotesIndexPages .Cfg.Layout.Numbers.NotesOnPage }}\hfill{}\hyperlink{Notes Index}{All notes}{{ end }}}
  \myMash[\myDailySpring]{\myNumDailyNotes}{\myNumDotWidthTwoThirds}
\end{minipage}
\par\pagebreak
