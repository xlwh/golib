module_state2��module_state�������档���������¹��ܣ�
1. ����֧��counter�⣬��֧������string���͵�state
2. ֧��counter slice��Ҳ���ǻ������counter�Ĳ�ֵ�����ڻ��
  һ��ʱ����counter�ı仯�����
  - �Ƽ���ʹ�÷�ʽ��
    ʹ��һ��������go routine, �����Եĵ���CounterSlice��Set()
3. ֧�ֱ�ƽ�����ṹ����Ӧ�Ĳ�λ������ṹjson���ת��


�ļ��б�
- counter.go: ʵ��Counters
  * Counters����Ϊmap[string]int64
- counter_slice.go��ʵ��CounterSlice
  * CounterSlice���ڻ����һ��ʱ���ڡ�Counters���������ռ�ı仯
- interval.go
  * ʵ����NextInterval()�������жϾ���һ����ʣ���ʱ��
- module_state2.go: ʵ��State
  * ����counter, �ַ������͵�state, �������͵�state
- counter_convert.go
    * �ṩ��Counters�ṹ���������Ĺ���
- counter_hier.go :
    * �ṩ��λ���counter�ṹ����ƽcounter�����counter��ת������
- counter_slice_hier.go :
    * �ṩ��ƽ��counter slice�ṹ����λ�counter slice�ṹ��json���ת��
- module_state2_hier.go :
    * �ṩ��ƽ��module״̬ͳ�ƽṹ����λ�module״̬ͳ�ƽṹ��json���ת��
