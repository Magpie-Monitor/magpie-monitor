import SVGIcon from 'components/SVGIcon/SVGIcon';
import './HiddenWebhook.scss';
import { useReducer } from 'react';

export interface HiddenWebhookProps {
  url: string;
}

const HiddenWebhook = ({ url }: HiddenWebhookProps) => {
  const [isHidden, toggle] = useReducer(
    (hiddenToChange) => !hiddenToChange,
    true,
  );

  const hiddenUrl: string = url
    .split('/')
    .map((segment, index, array) =>
      index === array.length - 1 ? '*'.repeat(segment.length) : segment,
    )
    .join('/');

  return (
    <div className="hidden-webhook">
      <div className='hidden-webhook__url'>{isHidden ? hiddenUrl : url}</div>
      <div onClick={toggle} className="hidden-webhook__button">
        <SVGIcon iconName={isHidden ? 'eye-icon' : 'closed-eye-icon'} />
      </div>
    </div>
  );
};

export default HiddenWebhook;
