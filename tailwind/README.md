# Install node/npm. Then install tailwind cli and init:

npm install -D tailwindcss
npx tailwindcss init

# rebuild static Tailwind CSS with

npx tailwindcss -i ./styles.css -o ../src/assets/styles.css --watch
