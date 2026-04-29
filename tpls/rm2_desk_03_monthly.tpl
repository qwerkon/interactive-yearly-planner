{{- $month := .Body.Month -}}
{{- $monthly := .Cfg.Layout.DeskMonthly -}}
\noindent
\begin{minipage}[t]{\dimexpr\linewidth-{{ $monthly.SidePanelWidth }}-4mm}
  \vspace{0pt}
  {\fontsize{ {{- $monthly.HeaderFont -}} }{ {{- $monthly.HeaderFont -}} }\selectfont\textbf{ {{- $month.Month.String -}} }}\hspace{2mm}{\large {{ .Body.HeadingMOS }} }
  
  \vspace{1mm}
  \myLineThick
\end{minipage}%
{{ if $monthly.ShowSidePanel -}}
\hfill%
\begin{minipage}[t]{ {{- $monthly.SidePanelWidth -}} }
  \vspace{0pt}\raggedleft
  {\fontsize{ {{- $monthly.SubHeaderFont -}} }{8}\selectfont {{ .Body.Breadcrumb }} }
  \vspace{1mm}

  {\fontsize{ {{- $monthly.LegendFont -}} }{6}\selectfont
  \textcolor{ {{- .Cfg.Layout.Colors.Saturday -}} }{Saturday}\hspace{1mm}
  \textcolor{ {{- .Cfg.Layout.Colors.PublicHoliday -}} }{Sunday/holiday}}
\end{minipage}
{{- end }}

\vspace{1mm}

{{ template "monthTabularV2.tpl" dict "Cfg" .Cfg "Month" $month "Today" nil "TableType" "tabularx" "Large" true }}

\vfill
\pagebreak
