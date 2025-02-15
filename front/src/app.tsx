import { defineComponent, onMounted, onUnmounted } from 'vue';
import Home from '@/views/home';
import { useUserStore } from '@/store';
export default defineComponent({
  setup() {
    // const router = useRouter();
    // 设置 rem
    const calcRem = () => {
      const doc = window.document;
      const docEl = doc.documentElement;
      const designWidth = 1580; // 默认设计图宽度
      const maxRate = 2560 / designWidth;
      const minRate = 1280 / designWidth;
      const clientWidth = docEl.getBoundingClientRect().width || window.innerWidth;
      const flexibleRem = Math.max(Math.min(clientWidth / designWidth, maxRate), minRate) * 100;
      docEl.style.fontSize = `${flexibleRem}px`;
    };
    const userStore = useUserStore();
    userStore.userInfo();
    onMounted(() => {
      calcRem();
      window.addEventListener('resize', calcRem, false);
    });
    onUnmounted(() => {
      window.removeEventListener('resize', calcRem, false);
    });
    return () => (
      <Home></Home>
    );
  },
});
