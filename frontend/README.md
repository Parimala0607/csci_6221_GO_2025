# Getting Started with Create React App

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Available Scripts

In the project directory, you can run:

### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in your browser.

The page will reload when you make changes.\
You may also see any lint errors in the console.

### `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can't go back!**

If you aren't satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you're on your own.

You don't have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn't feel obligated to use this feature. However we understand that this tool wouldn't be useful if you couldn't customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

This section has moved here: [https://facebook.github.io/create-react-app/docs/code-splitting](https://facebook.github.io/create-react-app/docs/code-splitting)

### Analyzing the Bundle Size

This section has moved here: [https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size](https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size)

### Making a Progressive Web App

This section has moved here: [https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app](https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app)

### Advanced Configuration

This section has moved here: [https://facebook.github.io/create-react-app/docs/advanced-configuration](https://facebook.github.io/create-react-app/docs/advanced-configuration)

### Deployment

This section has moved here: [https://facebook.github.io/create-react-app/docs/deployment](https://facebook.github.io/create-react-app/docs/deployment)

### `npm run build` fails to minify

This section has moved here: [https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify](https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify)

Required installations

1. VS-Code
2. Docker Desktop
3. GO lang 
4. Node js

VS-code extension:
1. Javascript and Typescript - microsoft
2. Dev Containers
3. Docker
4. ES7 React
5. ESLint
6. Go
7. React Developer Tools
8. SQLite Viewer
9. WSL


Debugging and Running:
React Frontend → http://localhost:3000
Dashboard API → http://localhost:8080/api/alerts & /logs
Blue Team → http://localhost:8081/defend
Red Team → http://localhost:8082/attack

docker-compose down -v
docker-compose build --no-cache
docker-compose up

Front-end:
npm install
npm audit fix --force
npm start
npm run build
cd frontend
  npm start


go run blueteam/main.go
go run redteam/main.go


Verify DB manually :

D:\go-work\GO_project\redblue-sim>docker volume inspect simdata --format '{{.Mountpoint}}'
'/var/lib/docker/volumes/simdata/_data'

D:\go-work\GO_project\redblue-sim>explorer $(docker volume inspect simdata --format '{{.Mountpoint}}')

D:\go-work\GO_project\redblue-sim>docker ps --filter "name=dashboard" --format "{{.ID}}"
21846f44a7d2

D:\go-work\GO_project\redblue-sim>docker exec -it 21846f44a7d2 sh
/app # apk add sqlite
OK: 10 MiB in 20 packages
/app # sqlite3 /data/sim.db
SQLite version 3.41.2 2023-03-22 11:56:21
Enter ".help" for usage hints.
sqlite> .tables
alerts  logs    users 