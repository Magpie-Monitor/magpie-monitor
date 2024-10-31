import { useReducer, useEffect, useRef } from 'react';
import './TagButton.scss';
import SVGIcon from 'components/SVGIcon/SVGIcon';

interface TagButtonProps {
  listItems: string[];
  chosenItem: string;
  onSelect: (item: string) => void;
}

interface State {
  isOpen: boolean;
  selectedOption: string;
}

type Action =
    | { type: 'TOGGLE' }
    | { type: 'SELECT_OPTION'; payload: string }
    | { type: 'CLOSE' };

const reducer = (state: State, action: Action): State => {
  switch (action.type) {
    case 'TOGGLE':
      return { ...state, isOpen: !state.isOpen };
    case 'SELECT_OPTION':
      return { isOpen: false, selectedOption: action.payload };
    case 'CLOSE':
      return { ...state, isOpen: false };
    default:
      return state;
  }
};

const TagButton = ({ listItems, chosenItem = listItems[0] }: TagButtonProps) => {
  const [{ isOpen, selectedOption }, dispatch] = useReducer(reducer, {
    isOpen: false,
    selectedOption: chosenItem,
  });

  const menuRef = useRef<HTMLUListElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
          menuRef.current &&
          !menuRef.current.contains(event.target as Node) &&
          !(event.target as HTMLElement).closest('.tag-button__toggle')
      ) {
        dispatch({ type: 'CLOSE' });
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleBlur = (e: React.FocusEvent<HTMLButtonElement>) => {
    if (!e.currentTarget.contains(e.relatedTarget as Node)) {
      dispatch({ type: 'CLOSE' });
    }
  };

  return (
      <div className="tag-button">
        <button
            className="tag-button__toggle"
            onClick={() => dispatch({ type: 'TOGGLE' })}
            onBlur={handleBlur}
            aria-haspopup="true"
            aria-expanded={isOpen}
        >
        <span className="tag-button__toggle__description">
          {selectedOption || '^'}
        </span>
          <SVGIcon iconName={isOpen ? 'reverse-drop-down-icon' : 'drop-down-icon'} />
        </button>
        {isOpen && (
            <ul className="tag-button__menu" ref={menuRef} role="menu">
              {listItems.map((item, index) => (
                  <li
                      key={index}
                      className="tag-button__menu__element"
                      onClick={() => dispatch({ type: 'SELECT_OPTION', payload: item })}
                      role="menuitem"
                      tabIndex={0}
                      onKeyDown={(e) => e.key === 'Enter' && dispatch({ type: 'SELECT_OPTION', payload: item })}
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