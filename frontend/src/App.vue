<script setup lang="ts">
import { RouterView, RouterLink } from 'vue-router'
</script>

<template>
  <header class="app-header">
    <RouterLink to="/" style="text-decoration: none; color: inherit;">
      <h1>Loup-Garou</h1>
    </RouterLink>
  </header>
  <main class="container">
    <RouterView />
  </main>
</template>
