# Tailwind local build

This code can be used to generate a light-weight CSS file with Tailwind styles that are actually used.

## Prerequisites

- node
- npm

Run `npm install`

## rebuild static Tailwind CSS with

```
npx tailwindcss -i ./styles.css -o ../src/assets/styles.css
```

## How this was set up

- npm install -D tailwindcss
- npx tailwindcss init
