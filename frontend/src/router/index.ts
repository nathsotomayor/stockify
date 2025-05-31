import { createRouter, createWebHistory } from 'vue-router'
import StockListView from '../views/StockListView.vue'
import StockDetailView from '../views/StockDetailView.vue'

const routes = [
    {
        path: '/',
        name: 'StockList',
        component: StockListView,
    },
    {
        path: '/stock/:ticker',
        name: 'StockDetail',
        component: StockDetailView,
        props: true,
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
})

export default router
