<script setup lang="ts">
import {
  ref,
  computed,
  provide,
} from 'vue';

import HostManage from '@/views/resource/resource-manage/children/manage/host-manage.vue';
import VpcManage from '@/views/resource/resource-manage/children/manage/vpc-manage.vue';
import SubnetManage from '@/views/resource/resource-manage/children/manage/subnet-manage.vue';
import SecurityManage from '@/views/resource/resource-manage/children/manage/security-manage.vue';
import DriveManage from '@/views/resource/resource-manage/children/manage/drive-manage.vue';
import IpManage from '@/views/resource/resource-manage/children/manage/ip-manage.vue';
import RoutingManage from '@/views/resource/resource-manage/children/manage/routing-manage.vue';
import ImageManage from '@/views/resource/resource-manage/children/manage/image-manage.vue';
import NetworkInterfaceManage from '@/views/resource/resource-manage/children/manage/network-interface-manage.vue';
import recyclebinManage from '@/views/resource/recyclebin-manager/recyclebin-manager.vue';
import { useVerify } from '@/hooks';
import useAdd from '@/views/resource/resource-manage/hooks/use-add';
import GcpAdd from '@/views/resource/resource-manage/children/add/gcp-add';
// forms
import EipForm from './forms/eip/index.vue';
import subnetForm from './forms/subnet/index.vue';
import securityForm from './forms/security/index.vue';

import {
  useRoute,
  useRouter,
} from 'vue-router';

import { useAccountStore } from '@/store/account';

const isShowSideSlider = ref(false);
const isShowGcpAdd = ref(false);
const componentRef = ref();
const securityType = ref('group');

// use hooks
const route = useRoute();
const router = useRouter();
const accountStore = useAccountStore();

const gcpTitle = ref<string>('新增');
const isAdd = ref(true);
const isLoading = ref(false);

provide('securityType', securityType);    // 将数据传入孙组件

// 组件map
const componentMap = {
  host: HostManage,
  vpc: VpcManage,
  subnet: SubnetManage,
  security: SecurityManage,
  drive: DriveManage,
  ip: IpManage,
  routing: RoutingManage,
  image: ImageManage,
  'network-interface': NetworkInterfaceManage,
  recyclebin: recyclebinManage,
};
const formMap = {
  ip: EipForm,
  subnet: subnetForm,
  security: securityForm,
};

const filter = ref({ op: 'and', rules: [] });

const renderComponent = computed(() => {
  return Object.keys(componentMap).reduce((acc, cur) => {
    if (route.path.includes(cur)) acc = componentMap[cur];
    return acc;
  }, {});
});

const renderForm = computed(() => {
  return Object.keys(formMap).reduce((acc, cur) => {
    if (route.path.includes(cur)) acc = formMap[cur];
    return acc;
  }, {});
});

const isResourcePage = computed(() => {   // 资源下没有业务ID
  return !accountStore.bizs;
});

const handleAdd = () => {
  if (renderComponent.value === DriveManage) {
    router.push({
      path: '/service/service-apply/disk',
    });
  } else if (renderComponent.value === HostManage) {
    router.push({
      path: '/service/service-apply/cvm',
    });
  } else if (renderComponent.value === VpcManage) {
    router.push({
      path: '/service/service-apply/vpc',
    });
  } else {
    if (securityType.value === 'gcp') {
      isShowGcpAdd.value = true;
      console.log('isShowGcpAdd.value', isShowGcpAdd.value);
    } else {
      isShowSideSlider.value = true;
    }
  }
};

const handleCancel = () => {
  isShowSideSlider.value = false;
};

// 新增成功 刷新列表
const handleSuccess = () => {
  handleCancel();
  componentRef.value.fetchComponentsData();
};

const handleSecrityType = (val: string) => {
  securityType.value = val;
  console.log(' securityType.value', securityType.value);
};

// 新增修改防火墙规则
const submit = async (data: any) => {
  const fetchType = 'vendors/gcp/firewalls/rules/create';
  const {
    addData,
    updateData,
  } = useAdd(
    fetchType,
    data,
    data?.id,
  );
  if (isAdd.value) {   // 新增
    addData();
  } else {
    await updateData();
  }
  isLoading.value = false;
};

const handleToPage = () => {
  const isHostManagePage = route.path.includes('/business/host');
  const isDriveManagePage = route.path.includes('/business/drive');
  let destination = '';
  if(isHostManagePage) destination = '/business/host/recyclebin/cvm';
  if(isDriveManagePage) destination = '/business/drive/recyclebin/disk';
  router.push({ path: destination });
}


// 权限hook
const {
  showPermissionDialog,
  handlePermissionConfirm,
  handlePermissionDialog,
  handleAuth,
  permissionParams,
  authVerifyData,
} = useVerify();
</script>

<template>
  <div>
    <section class="business-manage-wrapper">
      <bk-loading :loading="!accountStore.bizs">
        <component
          v-if="accountStore.bizs"
          ref="componentRef"
          :is="renderComponent"
          :filter="filter"
          :is-resource-page="isResourcePage"
          :auth-verify-data="authVerifyData"
          @auth="(val: string) => {
            handleAuth(val)
          }"
          @handleSecrityType="handleSecrityType"
        >
          <span @click="handleAuth('biz_iaas_resource_create')">
            <bk-button
              theme="primary" class="new-button"
              :disabled="!authVerifyData?.permissionAction?.biz_iaas_resource_create
                || securityType === 'gcp'" @click="handleAdd">
              {{renderComponent === DriveManage ||
                renderComponent === HostManage ||
                renderComponent === VpcManage ? '申请' : '新增'}}
            </bk-button>
          </span>

          <template #recycleHistory>
            <bk-button
              class="f-right"
              theme="primary"
              @click="handleToPage">
              {{ '回收记录' }}
            </bk-button>
          </template>
        </component>
      </bk-loading>
    </section>
    <bk-sideslider
      v-model:isShow="isShowSideSlider"
      width="800"
      title="新增"
      quick-close
    >
      <template #default>
        <component
          :is="renderForm" :filter="filter"
          @cancel="handleCancel" @success="handleSuccess"></component>
      </template>
    </bk-sideslider>
    <permission-dialog
      v-model:is-show="showPermissionDialog"
      :params="permissionParams"
      @cancel="handlePermissionDialog"
      @confirm="handlePermissionConfirm"
    ></permission-dialog>

    <gcp-add
      v-model:is-show="isShowGcpAdd"
      :gcp-title="gcpTitle"
      :is-add="isAdd"
      :loading="isLoading"
      :detail="{}"
      @submit="submit"></gcp-add>
  </div>
</template>

<style lang="scss" scoped>
.business-manage-wrapper {
  height: calc(100% - 20px);
  overflow-y: auto;
  background-color: #fff;
  padding: 20px;
}
.new-button {
  width: 100px;
}
</style>
