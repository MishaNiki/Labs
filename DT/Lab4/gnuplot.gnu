#! /usr/bin/gnuplot -persist

set term png enhanced font Hack 12 size 1280, 960
# ffffff 000000 size 1280,960


set output "./gr.png" #указываем путь и имя файла

set tmargin 1
set size 1.0,1.0
# set nokey
set mxtics 2
set mytics 2
set grid xtics ytics mxtics mytics


set xlabel "" font "Hack, 12"
set ylabel "" font "Hack , 12"

plot "./output" using 1:2 with linespoints pointtype 0 linecolor 3 pointsize 2 linewidth 2 