import './TagButton.scss';
import React, {
  useEffect,
  useRef,
  useReducer,
  ReactElement,
  useState,
} from 'react';
import SVGIcon from 'components/SVGIcon/SVGIcon.tsx';
import { debounce } from 'lib/debounce';

const ITEM_BUTTON_HEIGHT = 40;

interface TagButtonProps<T extends ReactElement | string | number> {
  listItems: T[];
  chosenItem: T;
  onSelect: (item: T) => void;
}

type Action<T extends ReactElement | string | number> =
  | { type: 'TOGGLE' }
  | { type: 'CLOSE' }
  | { type: 'SELECT_OPTION'; payload: T };

interface State<T extends ReactElement | string | number> {
  isOpen: boolean;
  selectedOption: T;
}

function reducer<T extends ReactElement | string | number>(
  state: State<T>,
  action: Action<T>,
): State<T> {
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

const TagButton = <T extends ReactElement | string | number>({
  listItems,
  chosenItem,
  onSelect,
}: TagButtonProps<T>): React.ReactNode => {
  const [{ isOpen, selectedOption }, dispatch] = useReducer(reducer, {
    isOpen: false,
    selectedOption: chosenItem,
  });

  const menuRef = useRef<HTMLUListElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const [itemPositions, setItemPositions] = useState<
    { top: number; left: number; width: number }[]
  >([]);

  const closeMenu = debounce(() => {
    dispatch({ type: 'CLOSE' });
  }, 10);

  const toggleMenu = () => {
    dispatch({ type: 'TOGGLE' });
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        menuRef.current &&
        !menuRef.current.contains(event.target as Node) &&
        !buttonRef.current?.contains(event.target as Node)
      ) {
        closeMenu();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  useEffect(() => {
    const updatePositions = () => {
      if (!buttonRef.current) return;

      const buttonRect = buttonRef.current.getBoundingClientRect();
      const positions = listItems.map((_, index) => ({
        top: buttonRect.bottom + index * ITEM_BUTTON_HEIGHT + window.scrollY,
        left: buttonRect.left + window.scrollX,
        width: buttonRect.width,
      }));
      setItemPositions(positions);
    };

    document.addEventListener('resize', closeMenu);
    document.addEventListener('scroll', closeMenu, true);
    updatePositions();

    return () => {
      document.removeEventListener('resize', closeMenu);
      document.removeEventListener('scroll', closeMenu);
    };
  }, [listItems, isOpen]);

  const handleSelect = (item: T) => {
    dispatch({ type: 'SELECT_OPTION', payload: item });
    onSelect(item);
    buttonRef.current?.focus();
  };

  return (
    <div className="tag-button">
      <button
        className="tag-button__toggle"
        ref={buttonRef}
        onClick={() => dispatch({ type: 'TOGGLE' })}
        onKeyDown={toggleMenu}
        aria-haspopup="listbox"
        aria-expanded={isOpen}
      >
        <span className="tag-button__toggle__description">
          {selectedOption}
        </span>
        <SVGIcon
          iconName={isOpen ? 'reverse-drop-down-icon' : 'drop-down-icon'}
        />
      </button>
      {isOpen && (
        <ul
          className="tag-button__menu"
          ref={menuRef}
          role="listbox"
          tabIndex={-1}
          style={{
            top: `${itemPositions.length > 0 ? itemPositions[0]?.top : 0}px`,
            left: `${itemPositions.length > 0 ? itemPositions[0]?.left : 0}px`,
          }}
        >
          {listItems.map((item, index) => (
            <li
              id={`tag-button-option-${index}`}
              key={index}
              className="tag-button__menu__element"
              onClick={() => handleSelect(item)}
              role="option"
              onKeyDown={(e) => {
                if (e.key === 'Enter') handleSelect(item);
              }}
              style={{
                top: `${itemPositions[index]?.top || 0}px`,
                left: `${itemPositions[index]?.left || 0}px`,
                width: `${itemPositions[index]?.width || 0}px`,
              }}
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
