import './ActionButton.scss';

export enum ActionButtonColor {
  GREEN,
  OLIVE,
  RED
}

export interface ActionButtonProps {
  onClick: () => void;
  description: string;
  color: ActionButtonColor;
}

const actionButtonColorToClass = {
  [ActionButtonColor.GREEN]: 'action-button--green',
  [ActionButtonColor.RED]: 'action-button--red',
  [ActionButtonColor.OLIVE]: 'action-button--olive',
};

const ActionButton = ({ onClick, description, color }: ActionButtonProps) => {
  return (
    <button onClick={onClick} className={actionButtonColorToClass[color]}>
      {description}
    </button>
  );
};

export default ActionButton;
