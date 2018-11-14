# EU4 Song List Generator

This command enables you to add your own song list to EU4 by generating `songs.txt` and `music.asset` from local `.ogg` files.

## Install

    go get -u github.com/tadayosi/eu4-songs-gen

## Usage

Running the command searches `.ogg` files in the current directory and generates `songs.txt` and `music.asset` there:

    eu4-songs-gen

You can also specify the directory where `.ogg` files are located and `songs.txt` and `music.asset` are generated there:

    eu4-songs-gen songs/

For more info, see the help:

    eu4-songs-gen -h
