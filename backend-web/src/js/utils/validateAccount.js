import { ElMessage } from 'element-plus'


export function validateContact (form) {
  const email = form.email?.trim()
  const phone = form.phone?.trim()

  if (!email && !phone) {
    // Please enter at least an email or a phone number
    ElMessage.warning('Either email or phone must be provided')
    return false
  }

  return true
}

const validateLabels = (rule, value, callback) => {
  if (!Array.isArray(value)) {
    callback();
    return;
  }

  const tooLong = value.find(label => label.length > 32);
  if (tooLong) {
    callback(new Error(`Label ${tooLong} exceeds 32 characters`));
    return;
  }

  if (value.length > 16) {
    callback(new Error('You can select up to 16 labels only'));
    return;
  }

  callback();
}

export const validateAccount = {
  firstname: [{ required: true, min: 1, max: 32, message: 'Please enter firstname', trigger: 'blur' }],
  lastname: [{ required: true, min: 1, max: 32, message: 'Please enter lastname', trigger: 'blur' }],
  email: [
    { type: 'email', message: 'invalid email format', trigger: 'blur' },
    { min: 5, max: 64, message: 'Email too long', trigger: 'blur' },
  ],
  phone: [{ min: 6, max: 20, trigger: 'blur' }],
  password: [{
    required: true, message: 'Please enter password',
    trigger: 'blur', min: 8, max: 32,
    pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,32}$/,
    // pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,32}$/,
    message: 'Password must be 8-32 chars with at least one uppercase, lowercase and number',
  }],
  level: [{ required: true, message: 'Please enter level', trigger: 'change' }],
  labels: [{ validator: validateLabels, trigger: 'change' }],
}
