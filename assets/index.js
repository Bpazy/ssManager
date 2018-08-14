const Foo = { template: '<div>foo</div>' };
const Bar = { template: '<div>bar</div>' };


const routes = [
    { path: '/foo', component: Foo },
    { path: '/bar', component: Bar }
];

const router = new VueRouter({
    routes // (缩写) 相当于 routes: routes
});

const app = new Vue({
    router
}).$mount('#app');
