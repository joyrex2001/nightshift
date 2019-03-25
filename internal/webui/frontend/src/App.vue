<template>
  <div id="app">
    <div id="nav">
          <h1>NIGHTSHIFT admin</h1>
          <router-link to="/scanners">Scanners</router-link> |
          <router-link to="/objects">Objects</router-link> |
          <router-link to="/about">About</router-link>
    </div>
    <router-view/>
  </div>
</template>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-size: 0.9em;
  color: #000000;
}

#nav {
  padding: 5px;
  background-image: linear-gradient(#000000, white);
}

#nav h1 {
  color: #fff;
  font-weight: bold;
  margin: 0px;
  text-shadow: 2px 2px #000000;
}

#nav a {
  font-weight: bold;
  color: #0060a0;
}

#nav a.router-link-exact-active {
  color: #0080ff;
}
</style>
