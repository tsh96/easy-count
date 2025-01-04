<template lang="pug">
  .landing-container
    h1 Welcome to Easy Count
    .button-container
      router-link.nav-button(to="/customer-records")
        | Customer Records
      router-link.nav-button(to="/personal-records")
        | Personal Records
  </template>

<style scoped>
.landing-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
}

.button-container {
  display: flex;
  gap: 20px;
  margin-top: 30px;
}

.nav-button {
  padding: 15px 30px;
  font-size: 1.2rem;
  text-decoration: none;
  color: white;
  background-color: #4CAF50;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.nav-button:hover {
  background-color: #45a049;
}
</style>