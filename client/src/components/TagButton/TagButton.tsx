import './TagButton.scss';
import React, { useEffect, useRef, useReducer } from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';

interface TagButtonProps<T> {
  listItems: T[];
  chosenItem: T;
  onSelect: (item: T) => void;
}

type Action<T> =
    | { type: 'TOGGLE' }
    | { type: 'CLOSE' }
    | { type: 'SELECT_OPTION'; payload: T };

interface State<T> {
  isOpen: boolean;
  selectedOption: T;
}

function reducer<T>(state: State<T>, action: Action<T>): State<T> {
  switch (action.type) {
    case 'TOGGLE':
      return { ...state, isOpen: !state.isOpen };
    case 'CLOSE':
      return { ...state, isOpen: false };
    case 'SELECT_OPTION':
      return { isOpen: false, selectedOption: action.payload };
    default:
      return state;
  }
}

// eslint-disable-next-line @typescript-eslint/explicit-function-return-type
const TagButton = <T,>({ listItems, chosenItem, onSelect }: TagButtonProps<T>) => {
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

  const handleSelect = (item: T) => {
    dispatch({ type: 'SELECT_OPTION', payload: item });
    onSelect(item);
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
          {selectedOption as string}
        </span>
          <SVGIcon iconName={isOpen ? 'reverse-drop-down-icon' : 'drop-down-icon'} />
        </button>
        {isOpen && (
            <ul className="tag-button__menu" ref={menuRef} role="menu">
              {listItems.map((item, index) => (
                  <li
                      key={index}
                      className="tag-button__menu__element"
                      onClick={() => handleSelect(item)}
                      role="menuitem"
                      tabIndex={0}
                      onKeyDown={(e) => e.key === 'Enter' && handleSelect(item)}
                  >
                    {item as string}
                  </li>
              ))}
            </ul>
        )}
      </div>
  );
};

export default TagButton;
