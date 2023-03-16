import {
  Table,
  Loading,
  Radio,
  Message,
} from 'bkui-vue';
import {
  defineComponent,
  h,
  ref,
  computed,
} from 'vue';
import {
  useI18n,
} from 'vue-i18n';
import StepDialog from '@/components/step-dialog/step-dialog';
import useQueryList  from '../../../hooks/use-query-list';
import useColumns from '../../../hooks/use-columns';
import {
  useResourceStore
} from '@/store/resource';

// 主机选硬盘挂载
export default defineComponent({
  components: {
    StepDialog,
  },

  props: {
    title: {
      type: String,
    },
    isShow: {
      type: Boolean,
    },
    detail: {
      type: Object,
    }
  },

  emits: ['update:isShow', 'success'],

  setup(props, { emit }) {
    const {
      t,
    } = useI18n();

    const deviceName = ref();
    const cachingType = ref();

    const cacheTypes = [
      'None',
      'ReadOnly',
      'ReadWrite'
    ]

    const rules = [
      {
        field: 'vendor',
        op: 'eq',
        value: props.detail.vendor,
      },
      {
        field: 'account_id',
        op: 'eq',
        value: props.detail.account_id,
      },
      {
        field: 'zone',
        op: 'eq',
        value: props.detail.zone,
      },
      {
        field: 'region',
        op: 'eq',
        value: props.detail.region,
      }
    ]

    if (props.detail.vendor === 'azure') {
      rules.push({
        field: 'resource_group_name',
        op: 'eq',
        value: props.detail.resource_group_name
      })
    }

    const {
      datas,
      pagination,
      isLoading,
      handlePageChange,
      handlePageSizeChange,
      handleSort,
    } = useQueryList(
      {
        filter: {
          op: 'and',
          rules,
        },
      },
      'disks'
    );

    const columns = useColumns('drive', true);

    const resourceStore = useResourceStore();

    const selection = ref<any>({});

    const isConfirmLoading = ref(false);

    const renderColumns = [
      {
        label: 'ID',
        field: 'id',
        render({ data }: any) {
          return h(
            Radio,
            {
              'model-value': selection.value.id,
              label: data.id,
              key: data.id,
              onChange() {
                selection.value = data;
              },
            }
          );
        },
      },
      ...columns
    ]

    const renderList = computed(() => {
      return datas.value.map((data) => !data.instance_id)
    })

    // 方法
    const handleClose = () => {
      emit('update:isShow', false);
    };

    const handleConfirm = () => {
      isConfirmLoading.value = true;
      const postData: any = {
        disk_id: selection.value.id,
        cvm_id: props.detail.id,
      }
      if (!selection.value.id) {
        Message({
          theme: 'error',
          message: '请先选择云硬盘'
        })
        return
      }
      if (props.detail.vendor === 'aws') {
        if (!deviceName.value) {
          Message({
            theme: 'error',
            message: '请先输入设备名称'
          })
          return
        }
        postData.device_name = deviceName.value
      }
      if (props.detail.vendor === 'azure') {
        if (!cachingType.value) {
          Message({
            theme: 'error',
            message: '请先选择缓存类型'
          })
          return
        }
        postData.caching_type = cachingType.value
      }
      resourceStore.attachDisk(postData).then(() => {
        emit('success');
        handleClose();
      }).catch((err: any) => {
        Message({
          theme: 'error',
          message: err.message || err
        })
      }).finally(() => {
        isConfirmLoading.value = false;
      })
    };

    return {
      deviceName,
      cachingType,
      cacheTypes,
      renderList,
      pagination,
      isLoading,
      renderColumns,
      isConfirmLoading,
      handlePageChange,
      handlePageSizeChange,
      handleSort,
      t,
      handleClose,
      handleConfirm,
    };
  },

  render() {
    const steps = [
      {
        isConfirmLoading: this.isConfirmLoading,
        component: () =>
          <Loading loading={this.isLoading}>
            {
              this.detail.vendor === 'aws'
              ? <>
                <span class="mr10">设备名称:</span>
                <bk-input v-model={this.deviceName} style="width: 200px;"></bk-input>
                </>
              : ''
            }
            {
              this.detail.vendor === 'azure'
              ? <>
                <span class="mr10">缓存类型:</span>
                <bk-select v-model={this.cachingType} style="width: 200px;display: inline-block;">
                  {
                    this.cacheTypes.map((type) => <bk-option
                      key={type}
                      value={type}
                      label={type}
                  />)
                  }
                </bk-select>
              </>
              : ''
            }
            <Table
              class="mt20"
              row-hover="auto"
              remote-pagination
              pagination={this.pagination}
              columns={this.renderColumns}
              data={this.renderList}
              onPageLimitChange={this.handlePageSizeChange}
              onPageValueChange={this.handlePageChange}
              onColumnSort={this.handleSort}
            />
          </Loading>
      },
    ];

    return <>
      <step-dialog
        title={this.t('挂载云硬盘')}
        isShow={this.isShow}
        steps={steps}
        onConfirm={this.handleConfirm}
        onCancel={this.handleClose}
      >
      </step-dialog>
    </>;
  },
});
