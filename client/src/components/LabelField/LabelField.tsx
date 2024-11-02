import './LabelField.scss';

interface LabelFieldParams {
  label: string;
  field: string;
}

const LabelField = ({ label, field }: LabelFieldParams) => {
  return (
    <div className="label-field">
      <div className="label-field__label">{label}</div>
      <div className="label-field__field">{field}</div>
    </div>
  );
};

export default LabelField;
