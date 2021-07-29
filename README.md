Usage: 

    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
        name: ecr-updater
    rules:
    - apiGroups: ["*"]
      resources: ["secrets"]
      verbs: ["*"]
      resourceNames: ["ecr-secret"]
    - apiGroups: ["*"] # this is reaquired as create doesn't work with resourceNames
      resources: ["secrets"]
      verbs: ["create"]
    - apiGroups: [""]
      resources: ["namespaces"]
      verbs: ["list"]
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
        name: ecr-updater
    subjects:
    - kind: ServiceAccount
      name: default
      namespace: ecr-updater
      roleRef:
      kind: ClusterRole
      name: ecr-updater
      apiGroup: rbac.authorization.k8s.io
    ---
    apiVersion: v1
    kind: Namespace
    metadata:
        name: ecr-updater
    ---
    apiVersion: v1
    kind: Secret
    metadata:
        name: aws-credentials
        namespace: ecr-updater
    stringData:
        AWS_ACCESS_KEY_ID: XXXXXXXXXXXXX
        AWS_SECRET_ACCESS_KEY: "YYYYYYYYYYYYYYYYYYYYYYYy"
    ---
    apiVersion: batch/v1beta1
    kind: CronJob
    metadata:
      name: ecr-secrets
      namespace: ecr-updater
    spec:
      schedule: "*/30 * * * *"
      concurrencyPolicy: Allow
      failedJobsHistoryLimit: 1
      successfulJobsHistoryLimit: 1
      jobTemplate:
        spec:
          backoffLimit: 0
          template:
            spec:
              containers:
              - name: hello
                image: bumbacea/ecr-k8s-updater:0.0.2
                imagePullPolicy: IfNotPresent
                env:
                - name: SECRET_NAME
                  value: ecr-secret
                - name: AWS_REGION
                  value: eu-central-1
                - name: AWS_SECRET_ACCESS_KEY
                  valueFrom:
                    secretKeyRef:
                      name: aws-credentials
                      key: AWS_SECRET_ACCESS_KEY
                - name: AWS_ACCESS_KEY_ID
                  valueFrom:
                    secretKeyRef:
                      name: aws-credentials
                      key: AWS_ACCESS_KEY_ID
              restartPolicy: Never


    
