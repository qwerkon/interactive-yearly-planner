{{- $days := .Body.Week.Days -}}
{{- $month := index .Body.Week.Months 0 -}}
\noindent
\begin{minipage}[t]{\dimexpr\linewidth-4.55cm}
  \vspace{1mm}
  {\fontsize{13}{13}\selectfont\textbf{ {{- range $i, $m := .Body.Week.Months -}}{{if $i}} / {{end}}{{ $m.Month.String }}{{- end -}} }}\hspace{2mm}{\large {{ .Body.Week.Target }}}
  \vspace{1mm}

  \myLineThick
\end{minipage}%
\hfill%
\begin{minipage}[t]{4.25cm}
  \raggedleft
  \tiny
  {{ template "monthTabularV2.tpl" dict "Cfg" .Cfg "Month" $month "Today" nil "TableType" "tabularx" "Large" false }}
\end{minipage}

\vspace{1mm}

\setlength{\tabcolsep}{0pt}%
\renewcommand{\arraystretch}{1}%
\providecommand{\deskDayHeader}[3]{%
  \begin{minipage}[t][13mm][t]{\dimexpr\linewidth-2mm}%
    \vspace{1mm}%
    \hspace*{1mm}{\fontsize{18}{18}\selectfont\textbf{#1}}\\[-.7mm]
    \hspace*{1mm}{\fontsize{6.5}{7}\selectfont #2}%
    \if\relax\detokenize{#3}\relax\else\\[-.4mm]
      \hspace*{1mm}{\fontsize{4.3}{4.7}\selectfont #3}%
    \fi%
  \end{minipage}%
}
\providecommand{\deskHourLine}[1]{%
  \vbox to 4.25mm{%
    \vfill%
    \hbox to \linewidth{%
      \makebox[7mm][l]{\fontsize{5.5}{6}\selectfont\textcolor{gray}{#1}}%
      \textcolor{gray!45}{\leaders\hrule height \myLenLineThicknessDefault\hfill}%
    }%
    \vfill%
  }%
}
\providecommand{\deskDayColumn}[2]{%
  \hspace*{1mm}%
  \begin{minipage}[t]{\dimexpr\linewidth-2mm}%
    \vspace*{1mm}%
    #1%
    \deskHourLine{06:00}%
    \deskHourLine{07:00}%
    \deskHourLine{08:00}%
    \deskHourLine{09:00}%
    \deskHourLine{10:00}%
    \deskHourLine{11:00}%
    \deskHourLine{12:00}%
    \deskHourLine{13:00}%
    \deskHourLine{14:00}%
    \deskHourLine{15:00}%
    \deskHourLine{16:00}%
    \deskHourLine{17:00}%
    \deskHourLine{18:00}%
    \deskHourLine{19:00}%
    \deskHourLine{20:00}%
    \deskHourLine{21:00}%
    \deskHourLine{22:00}%
    #2%
    \vspace*{1mm}%
  \end{minipage}%
  \hspace*{1mm}%
}

\begin{tabularx}{\linewidth}{*{7}{|>{\raggedright\arraybackslash}X}|}
\hline
{{- range $i, $day := $days }}
  {{- $color := $day.HeaderColor $.Cfg -}}
  {{- $holidayName := $day.HolidayLocalName $.Cfg -}}
  {{- if eq $i 6 -}}
    \textcolor{red!70!black}{\deskDayHeader{ {{- $day.Time.Day -}} }{ {{- $day.Time.Weekday.String -}} }{ {{- $holidayName -}} }}
  {{- else if eq $i 5 -}}
    \textcolor{gray!60!black}{\deskDayHeader{ {{- $day.Time.Day -}} }{ {{- $day.Time.Weekday.String -}} }{ {{- $holidayName -}} }}
  {{- else if $color -}}
    \textcolor{ {{- $color -}} }{\deskDayHeader{ {{- $day.Time.Day -}} }{ {{- $day.Time.Weekday.String -}} }{ {{- $holidayName -}} }}
  {{- else -}}
    \deskDayHeader{ {{- $day.Time.Day -}} }{ {{- $day.Time.Weekday.String -}} }{ {{- $holidayName -}} }
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
