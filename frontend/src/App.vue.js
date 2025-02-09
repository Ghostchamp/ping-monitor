import { defineComponent, onMounted, ref } from 'vue';
import axios from 'axios';
export default defineComponent({
    name: 'App',
    setup() {
        const pingResults = ref([]);
        const fetchData = async () => {
            try {
                const response = await axios.get('/api/stats');
                pingResults.value = response.data;
            }
            catch (error) {
                console.error('Error fetching data:', error);
            }
        };
        onMounted(() => {
            fetchData();
            // Обновляем данные каждые 10 секунд
            setInterval(fetchData, 10000);
        });
        return { pingResults };
    }
});
; /* PartiallyEnd: #3632/script.vue */
function __VLS_template() {
    const __VLS_ctx = {};
    let __VLS_components;
    let __VLS_directives;
    // CSS variable injection 
    // CSS variable injection end 
    __VLS_elementAsFunction(__VLS_intrinsicElements.div, __VLS_intrinsicElements.div)({
        ...{ style: ({}) },
    });
    const __VLS_0 = {}.ElTable;
    /** @type { [typeof __VLS_components.ElTable, typeof __VLS_components.elTable, typeof __VLS_components.ElTable, typeof __VLS_components.elTable, ] } */ ;
    // @ts-ignore
    const __VLS_1 = __VLS_asFunctionalComponent(__VLS_0, new __VLS_0({
        data: ((__VLS_ctx.pingResults)),
        ...{ style: ({}) },
    }));
    const __VLS_2 = __VLS_1({
        data: ((__VLS_ctx.pingResults)),
        ...{ style: ({}) },
    }, ...__VLS_functionalComponentArgsRest(__VLS_1));
    const __VLS_6 = {}.ElTableColumn;
    /** @type { [typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ] } */ ;
    // @ts-ignore
    const __VLS_7 = __VLS_asFunctionalComponent(__VLS_6, new __VLS_6({
        prop: ("ip"),
        label: ("IP Address"),
    }));
    const __VLS_8 = __VLS_7({
        prop: ("ip"),
        label: ("IP Address"),
    }, ...__VLS_functionalComponentArgsRest(__VLS_7));
    const __VLS_12 = {}.ElTableColumn;
    /** @type { [typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ] } */ ;
    // @ts-ignore
    const __VLS_13 = __VLS_asFunctionalComponent(__VLS_12, new __VLS_12({
        prop: ("ping_time"),
        label: ("Ping Time (ms)"),
    }));
    const __VLS_14 = __VLS_13({
        prop: ("ping_time"),
        label: ("Ping Time (ms)"),
    }, ...__VLS_functionalComponentArgsRest(__VLS_13));
    const __VLS_18 = {}.ElTableColumn;
    /** @type { [typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, typeof __VLS_components.ElTableColumn, typeof __VLS_components.elTableColumn, ] } */ ;
    // @ts-ignore
    const __VLS_19 = __VLS_asFunctionalComponent(__VLS_18, new __VLS_18({
        prop: ("last_success_at"),
        label: ("Last Successful Ping"),
    }));
    const __VLS_20 = __VLS_19({
        prop: ("last_success_at"),
        label: ("Last Successful Ping"),
    }, ...__VLS_functionalComponentArgsRest(__VLS_19));
    __VLS_5.slots.default;
    var __VLS_5;
    var __VLS_slots;
    var $slots;
    let __VLS_inheritedAttrs;
    var $attrs;
    const __VLS_refs = {};
    var $refs;
    var $el;
    return {
        attrs: {},
        slots: __VLS_slots,
        refs: $refs,
        rootEl: $el,
    };
}
;
let __VLS_self;
//# sourceMappingURL=App.vue.js.map