
import http from '@/http';
import { nextTick, reactive, ref, watch } from 'vue';
import { VendorEnum } from '@/common/constant';
import type { Cond } from './use-condtion';
import { Message } from 'bkui-vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

const { BK_HCM_AJAX_URL_PREFIX } = window.PROJECT_CONFIG;

export interface IDiskOption {
  disk_type: string;
  disk_size_gb: number;
  disk_count?: number;
  disk_name?: string;
  mode?: string;
  auto_delete?: boolean;
}

export interface IPurchaseDuration {
  count: number;
  unit: 'm' | 'y';
}

export interface ICvmBaseData {
  zone: string | string[];
  name: string;
  instance_type: string;
  cloud_image_id: string;
  cloud_vpc_id: string;
  cloud_subnet_id: string;
  public_ip_assigned?: boolean;
  cloud_security_group_ids?: string | string[];
  system_disk: IDiskOption;
  data_disk: IDiskOption[];
  username?: string;
  password: string;
  confirmed_password: string;
  instance_charge_type?: 'PREPAID' | 'POSTPAID_BY_HOUR' | 'SPOTPAID' | 'prePaid' | 'postPaid';
  instance_charge_paid_period?: number;
  auto_renew?: boolean;
  required_count: number;
  memo: string;
}
export interface ICvmFormData extends ICvmBaseData {
  purchase_duration: IPurchaseDuration;
}
export interface ICvmSaveData extends ICvmBaseData {
  bk_biz_id: number;
  account_id: string;
  region: string;
  vendor?: string;
  resource_group_name?: string;
}

export const getDataDiskDefaults = () => ({
  disk_type: '',
  disk_size_gb: 50,
  disk_count: 1,
});

export const getGcpDataDiskDefaults = (): IDiskOption => ({
  disk_type: '',
  disk_size_gb: 50,
  disk_count: 1,
  disk_name: '',
  mode: 'READ_WRITE',
  auto_delete: false,
});

export default (cond: Cond) => {
  const { t } = useI18n();
  const router = useRouter();

  const vendorDiffFormData = (vendor: string) => {
    const diff = {
      [VendorEnum.TCLOUD]: {
        instance_charge_type: 'PREPAID',
        public_ip_assigned: false,
        cloud_security_group_ids: [] as string[],
        auto_renew: false,
      },
      [VendorEnum.AWS]: {
        public_ip_assigned: false,
        cloud_security_group_ids: [] as string[],
      },
      [VendorEnum.AZURE]: {
        username: '',
        cloud_security_group_ids: '',
      },
      [VendorEnum.GCP]: {
        data_disk: [] as string[],
      },
      [VendorEnum.HUAWEI]: {
        public_ip_assigned: false,
        instance_charge_type: 'prePaid',
        cloud_security_group_ids: [] as string[],
        auto_renew: false,
      },
    };
    return diff[vendor] || {};
  };

  const defaultFormData = (vendor: string) => {
    const base: ICvmFormData = {
      zone: [],
      name: '',
      instance_type: '',
      cloud_image_id: '',
      cloud_vpc_id: '',
      cloud_subnet_id: '',
      system_disk: {
        disk_type: '',
        disk_size_gb: 50,
      },
      data_disk: [],
      password: '',
      confirmed_password: '',
      purchase_duration: {
        count: 1,
        unit: 'm',
      },
      required_count: 1,
      memo: '',
    };

    return {
      ...base,
      ...vendorDiffFormData(vendor),
    };
  };

  const formData = reactive<ICvmFormData>(defaultFormData(cond.vendor));
  const formRef = ref(null);

  const resetFormData = () => {
    const keys = [
      'zone',
      'instance_type',
      'cloud_image_id',
      'cloud_vpc_id',
      'cloud_subnet_id',
      'cloud_security_group_ids',
      'data_disk',
      'public_ip_assigned',
      'instance_charge_type',
      'auto_renew',
    ];
    keys.forEach(key => resetFormItemData(key));
  };

  const resetFormItemData = (key: string) => {
    const defaultData: ICvmFormData = defaultFormData(cond.vendor);
    formData[key] = defaultData[key];
  };

  watch(cond, (cond) => {
    resetFormData();

    nextTick(() => {
      formRef.value.clearValidate();
    });
  });

  const getSaveData = () => {
    const { purchase_duration, public_ip_assigned, ...saveFormData } = formData
    const saveData: ICvmSaveData = {
      ...saveFormData,
      bk_biz_id: cond.bizId,
      account_id: cond.cloudAccountId,
      region: cond.region,
      zone: formData?.zone?.[0],
    };

    if (cond.vendor === VendorEnum.TCLOUD) {
      saveData.public_ip_assigned = public_ip_assigned
      saveData.instance_charge_paid_period = purchase_duration.count * (purchase_duration.unit === 'y' ? 12 : 1)
    }

    if (cond.vendor === VendorEnum.HUAWEI) {
      saveData.public_ip_assigned = public_ip_assigned
      saveData.instance_charge_paid_period = purchase_duration.count * (purchase_duration.unit === 'y' ? 12 : 1)
    }

    if (cond.vendor === VendorEnum.AWS) {
      saveData.public_ip_assigned = public_ip_assigned
    }

    if (cond.vendor === VendorEnum.AZURE) {
      saveData.resource_group_name = cond.resourceGroup;
      saveData.cloud_security_group_ids = [saveFormData.cloud_security_group_ids as string]
    }
    return saveData;
  };

  const submitting = ref(false);
  const handleFormSubmit = async () => {
    await formRef.value.validate();
    const saveData = getSaveData();
    // console.log(saveData, '-----saveData');
    try {
      submitting.value = true;
      await http.post(`${BK_HCM_AJAX_URL_PREFIX}/api/v1/cloud/vendors/${cond.vendor}/applications/types/create_cvm`, saveData);

      Message({
        theme: 'success',
        message: t('提交成功'),
      });

      router.push({
        path: '/service/my-apply',
      });
    } catch (err) {
      console.error(err);
    } finally {
      submitting.value = false;
    }
  };

  return {
    formRef,
    formData,
    resetFormData,
    resetFormItemData,
    getSaveData,
    submitting,
    handleFormSubmit,
  };
};
