import { useReducer, useState } from 'react';
import './TagButton.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';

export interface TagButtonProps {
  listItems: Array<string>;
  chosenItem: string;
}

const TagButton = ({
  listItems,
  chosenItem = listItems[0],
}: TagButtonProps) => {
  const [isOpen, toggle] = useReducer(
    (isOpenToChange) => !isOpenToChange,
    false,
  );
  const [selectedOption, setSelectedOption] = useState(chosenItem);

  const handleOptionClick = (item: string) => {
    setSelectedOption(item);
    toggle();
  };

  return (
    <div className="tag-button">
      <div className="tag-button__toggle" onClick={toggle}>
        <div className="tag-button__toggle__description">{selectedOption}</div>
        <SVGIcon
          iconName={isOpen ? 'reverse-drop-down-icon' : 'drop-down-icon'}
        />
      </div>
      {isOpen && (
        <ul className="tag-button__menu">
          {listItems.map((item, index) => (
            <li
              key={index}
              className="tag-button__menu__element"
              onClick={() => handleOptionClick(item)}
            >
              {item}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default TagButton;
