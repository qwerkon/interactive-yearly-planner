{{- $days := .Body.Week.Days -}}
{{- $month := index .Body.Week.Months 0 -}}
{{- $desk := .Cfg.Layout.DeskWeekly -}}
\noindent
{{ if $desk.ShowMiniCalendar -}}
\begin{minipage}[t]{\dimexpr\linewidth-{{ $desk.MiniCalendarWidth }}-3mm}
{{- else -}}
\begin{minipage}[t]{\linewidth}
{{- end }}
  \vspace{1mm}
  {\fontsize{13}{13}\selectfont\textbf{ {{- range $i, $m := .Body.Week.Months -}}{{if $i}} / {{end}}{{ $m.Month.String }}{{- end -}} }}\hspace{2mm}{\large {{ .Body.Week.Target }}}
  \vspace{1mm}

  \myLineThick
\end{minipage}%
{{ if $desk.ShowMiniCalendar -}}
\hfill%
\begin{minipage}[t]{ {{- $desk.MiniCalendarWidth -}} }
  \raggedleft
  \tiny
  {{ template "monthTabularV2.tpl" dict "Cfg" .Cfg "Month" $month "Today" nil "TableType" "tabularx" "Large" false }}
\end{minipage}
{{- end }}

\vspace{1mm}

\setlength{\tabcolsep}{0pt}%
\renewcommand{\arraystretch}{1}%
\providecommand{\deskDayHeader}[4]{%
  \begin{minipage}[t][{{ $desk.HeaderHeight }}][t]{\dimexpr\linewidth-{{ $desk.HeaderPadding }}-{{ $desk.HeaderPadding }}\relax}%
    \vspace{ {{- $desk.HeaderPadding -}} }%
    \hspace*{ {{- $desk.HeaderPadding -}} }{\fontsize{ {{- $desk.DayNumberFont -}} }{ {{- $desk.DayNumberFont -}} }\selectfont\textbf{#1}}\hfill{\fontsize{ {{- $desk.HolidayNameFont -}} }{4.7}\selectfont #4}\\[-.7mm]
    \hspace*{ {{- $desk.HeaderPadding -}} }{\fontsize{ {{- $desk.WeekdayFont -}} }{7}\selectfont #2}%
    \if\relax\detokenize{#3}\relax\else\\[-.4mm]
      \hspace*{ {{- $desk.HeaderPadding -}} }{\fontsize{ {{- $desk.HolidayNameFont -}} }{4.7}\selectfont #3}%
    \fi%
  \end{minipage}%
}
\providecommand{\deskHourLine}[1]{%
  \vbox to {{ $desk.HourLineHeight }}{%
    \vfill%
    \hbox to \linewidth{%
      \makebox[7mm][l]{\fontsize{ {{- $desk.HourFont -}} }{6}\selectfont\textcolor{gray}{#1}}%
      \textcolor{gray!45}{\leaders\hrule height \myLenLineThicknessDefault\hfill}%
    }%
    \vfill%
  }%
}
\providecommand{\deskDayColumn}[2]{%
  \hspace*{ {{- $desk.ColumnPadding -}} }%
  \begin{minipage}[t]{\dimexpr\linewidth-{{ $desk.ColumnPadding }}-{{ $desk.ColumnPadding }}\relax}%
    \vspace*{ {{- $desk.ColumnPadding -}} }%
    #1%
    {{- range $hour := $desk.Hours }}
    \deskHourLine{ {{- $hour -}} }%
    {{- end }}
    #2%
    \vspace*{ {{- $desk.ColumnPadding -}} }%
  \end{minipage}%
  \hspace*{ {{- $desk.ColumnPadding -}} }%
}
\providecommand{\deskHolidayLegend}{%
  {{- if $desk.ShowHolidayLegend }}
  {\fontsize{5}{6}\selectfont\textcolor{ {{- $.Cfg.Layout.Colors.PublicHoliday -}} }{ {{- $desk.HolidayMarker -}}~holiday}\hspace{2mm}\textcolor{ {{- $.Cfg.Layout.Colors.Event -}} }{ {{- $desk.EventMarker -}}~event}}\vspace{.6mm}
  {{- end }}
}

\deskHolidayLegend

\begin{tabularx}{\linewidth}{*{7}{|>{\raggedright\arraybackslash}X}|}
\hline
{{- range $i, $day := $days }}
  {{- $color := $day.HeaderColor $.Cfg -}}
  {{- $holidayName := $day.HolidayLocalName $.Cfg -}}
  {{- $marker := $day.HolidayMarker $.Cfg -}}
  {{- if eq $i 6 -}}
    \textcolor{ {{- $.Cfg.Layout.Colors.Sunday -}} }{\deskDayHeader{ {{- $day.Time.Day -}} }{ {{- $day.Time.Weekday.String -}} }{ {{- $holidayName -}} }{ {{- $marker -}} }}
  {{- else if eq $i 5 -}}
    \textcolor{ {{- $.Cfg.Layout.Colors.Saturday -}} }{\deskDayHeader{ {{- $day.Time.Day -}} }{ {{- $day.Time.Weekday.String -}} }{ {{- $holidayName -}} }{ {{- $marker -}} }}
  {{- else if $color -}}
    \textcolor{ {{- $color -}} }{\deskDayHeader{ {{- $day.Time.Day -}} }{ {{- $day.Time.Weekday.String -}} }{ {{- $holidayName -}} }{ {{- $marker -}} }}
  {{- else -}}
    \deskDayHeader{ {{- $day.Time.Day -}} }{ {{- $day.Time.Weekday.String -}} }{ {{- $holidayName -}} }{ {{- $marker -}} }
  {{- end -}}
  {{- if eq $i 6 }}\\ \hline{{ else }} & {{ end -}}
{{- end }}
{{- range $i, $day := $days }}
  {{- if eq $i 6 -}}
    \deskDayColumn{}{}
  {{- else if eq $i 5 -}}
    \deskDayColumn{}{}
  {{- else -}}
    \deskDayColumn{}{}
  {{- end -}}
  {{- if eq $i 6 }}\\ \hline{{ else }} & {{ end -}}
{{- end }}
\end{tabularx}

\pagebreak
