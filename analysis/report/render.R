#!/usr/bin/env Rscript

# This script is intended to trigger the actual report generation using a
# R markdown file as a template.

# Install missing packages.
list.of.packages <- c("rmarkdown", "tidyverse", "lubridate")
new.packages <- list.of.packages[!(list.of.packages %in% installed.packages()[,"Package"])]
if(length(new.packages)) install.packages(new.packages)

library("rmarkdown")

# test that there are exactly 3 arguments
#  - the script to run
#  - the data source file
#  - the output directory
args = commandArgs(trailingOnly=TRUE)
if (length(args)!=4) {
  stop("Script requires exactly three parameters: <template> <data> <outputdir> <outputfile>", call.=FALSE)
}

template = args[1]
data = args[2]
outputdir = args[3]
outputfile = args[4]

rmarkdown::render(
    template,
    params = list(datafile=data),
    output_dir = outputdir,
    output_file = outputfile,
)