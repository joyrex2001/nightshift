<template>
  <div id="app">
    <div id="nav">
          <img alt="Nightshift" src="./assets/nightshift-logo.png" style="float:left; padding:0.6em;">
          <h1>Nightshift</h1>
          <router-link to="/scanners">Scanners</router-link> |
          <router-link to="/objects">Objects</router-link> |
          <router-link to="/about">About</router-link>
          <hr>
    </div>
    <router-view/>
  </div>
</template>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #000000;
}
#nav {
  padding: 2px;
}

#nav a {
  font-weight: bold;
  color: #004080;
}

#nav a.router-link-exact-active {
  color: #0080ff;
}
</style>
