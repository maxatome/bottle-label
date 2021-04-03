# bottle-label

## Installation

```
go install
```

## Usage

```
bottle-label -font georgia.json -output wine.pdf \
    "3 x Domaine Vayssette\nCuvée Léa - 2011\nGaillac - 9,20 -> 2016/9" \
    "6 x Château l'Hospitalet\nLa Clape - 2015\nLanguedoc - 7,85" \
    "Romain Duvernay\nVisan - 2016\nCôtes du Rhône V. - 5,30"
```

![bottle-label example](example.png?raw=true "bottle-label")

The `georgia.json` has to be generated with `makefont` utility of
`gofpdf` package, see https://github.com/jung-kurt/gofpdf#nonstandard-fonts

By default, `bottle-label` use helvetica font.

`-output` is optional. Default output file is `output.pdf`.
