<template>
  <div id="app-container" class="min-h-screen bg-slate-100 text-slate-800 flex flex-col">
    <header class="bg-white shadow-md sticky top-0 z-50">
      <nav class="container mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <router-link to="/" class="text-2xl font-bold text-indigo-600 hover:text-indigo-700 transition-colors">
          <span class="text-sky-500">Stockify</span>
        </router-link>
      </nav>
    </header>

    <main class="w-full flex-grow p-4 sm:p-6 lg:p-8">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <footer class="text-center py-6 border-t border-slate-200 bg-white">
      <p class="text-sm text-slate-500">
        &copy; {{ new Date().getFullYear() }} Nathaly Sotomayor Ampudia · Stockify
      </p>
    </footer>
  </div>
</template>

<style>
body {
  font-family: "Inter", sans-serif;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}
::-webkit-scrollbar-track {
  background: #f8fafc;
  border-radius: 10px;
}
::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 10px;
}
::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>
