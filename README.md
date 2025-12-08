# Test HTTP Server

Testing out running a Go HTTP server on Railway through their "Generate from GitHub Repo" option

Deploying is super easy, even updating the permissions to include a new repo is pretty easy. One of the options in Railway is "Configure GitHub App", clicking it navigates you to the permissions page in GitHub where you can add more repos to provide visibility for. Once updated, refreshing the page will show any newly added repos

Over 2 days the PostgreSQL, DB Init, and HTTP server instances cost 6 cents, over a month that'd be ~$1.86. This is all just services idling so cost could be much more, but for projects that'll be mostly at idle or low traffice Railway could be a solid option as deploying the various services is pretty easy
