# Gecko

A command line tool that makes it easier to style text, with great speed!

## Table of Contents
1. [Using Gecko](#using-gecko)
   1. [Backgrounds](#backgrounds)
   2. [Different Colors](#different-colors)
   3. [Styles](#styles)
   4. [Escaping Tags](#escaping-markup-tags)
2. [How to Install](#how-to-install)
3. [What's with the name?](#why-is-the-app-called-gecko-and-the-repository-called-chameleon)

## Using Gecko
To use Gecko, simply use the command with text.
```bash
$ gecko "Hello, World!"
```

But, this doesn't print with any color.
To use color, you need to add a markup tag. Like So:
```bash
$ gecko "[cyan]This text is cyan![/]"
```
Gecko follows the same markup tagging system as [Spectre.Console](https://github.com/spectreconsole/spectre.console/tree/main). It even uses the same color names, just with a couple small differences.

Unlike [Spectre.Console](https://github.com/spectreconsole/spectre.console/tree/main), Gecko doesn't read markup tags like a stack (doesn't require a trailing `[/]` to specify the end of each color segment). You can just switch colors on the fly, and the `[/]` tag has been changed to the reset color tag. It can be used anywhere in the input string, and it will reset the console to it's default. 

It's recommended that you still use the `[/]` tag, at least at the end of the line to prevent runaway colors from leaking into the users prompt, or terminals next output. An even better
solution is to add `\033[0m` to the front of your `PS1` environment variable to ensure nothing can mess with the next input prompt.

> [!IMPORTANT]  
> Anything that is found to be a tag will be parsed. If the data inside cannot be mapped
> to a color or style, then the data is ignored, effectively making them weirdly placed comments.

### Backgrounds

Gecko also implements an option to color the background of text. All you need to do is add the keyword `on` after the first color in a markup tag, then the color of the background color. Like so:

```bash
$ gecko "[cyan on white]This should look great![/]"
```

### Different Colors
Don't like the preset colors and want to use your own? Well you can use both hex an RGB values!
Gecko will adjust the color to the best of it's ability based on what the users terminal supports. You can use them like this:

```bash
$ gecko "[#7afb42 on rgb(0,10,10)]These colors make the console look hackery[/]"
```

### Styles

Following the ANSI chart of styles (and again, [Spectre.Console](https://github.com/spectreconsole/spectre.console/tree/main)), you can use styles with colors. 
It doesn't matter what the order of the words are either, just as long as the word is known by Gecko.

An example of styles being used with Gecko:

```bash
$ gecko "[blinking red]WARNING: I HAD TO LEARN ANOTHER PROGRAMMING LANGUAGE TO MAKE THIS, PLEASE USE IT[/]"
```

Or as mentioned above, it *any order*:

```bash
$ gecko "[underlined cyan blinking on black]wOW, THIs Is soO COoOl![/]"
```

### Escaping Markup Tags

To escape a markup tag, just add a second open bracket to the start of the tag:

```bash
$ gecko "[[this isn't going to be parsed for colors :D]"
# output: [this isn't going to be parsed for colors :D]
```

# How to Install

## Releases

Go to [Releases](releases) and download the most recent version.

## Building Yourself

Install the GoLang [compiler from The official site](https://go.dev/dl/), or install it with

```bash
$ sudo snap install go --classic
```

Clone the repo.

```bash
$ git clone "https://github.com/ScripturaOpus/ChameleonTerminal.git"
```

Navigate to `./ChameleonTerminal/`, then run `make`.
Gecko will be built into the `build` folder.

## Why is the app called "Gecko" and the repository called "Chameleon"?

Because the original name for Gecko was "Chameleon", but I couldn't find a shorter version of that name. 
So I changed it to Gecko. It might not be a color changing animal, 
but it's better than having to type out "chameleon" every time you wanted to use the command without having to set an alias for it.


# Sponsors

Why are you looking down here?

Aww, you really thought this was good enough that someone would give me money.
That's really nice.