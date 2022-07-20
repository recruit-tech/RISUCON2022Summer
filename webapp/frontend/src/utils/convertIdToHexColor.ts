import {
  blue,
  blueGrey,
  brown,
  cyan,
  deepOrange,
  deepPurple,
  green,
  grey,
  indigo,
  lightBlue,
  lightGreen,
  lime,
  pink,
  purple,
  red,
  teal,
} from "@mui/material/colors";

/**
 * colors in @mui/material/colors that meet the "WCAG 2.0 level AA" requirements
 * @link https://webaim.org/resources/contrastchecker/
 */
const COLOR_LIST = [
  blue["700"],
  blue["800"],
  blue["900"],
  blue["A700"],
  blueGrey["600"],
  blueGrey["700"],
  blueGrey["800"],
  blueGrey["900"],
  blueGrey["A700"],
  brown["400"],
  brown["500"],
  brown["600"],
  brown["700"],
  brown["800"],
  brown["900"],
  brown["A400"],
  brown["A700"],
  cyan["800"],
  cyan["900"],
  deepOrange["900"],
  deepOrange["A700"],
  deepPurple["400"],
  deepPurple["500"],
  deepPurple["600"],
  deepPurple["700"],
  deepPurple["800"],
  deepPurple["900"],
  deepPurple["A200"],
  deepPurple["A400"],
  deepPurple["A700"],
  green["800"],
  green["900"],
  grey["600"],
  grey["700"],
  grey["800"],
  grey["900"],
  grey["A700"],
  indigo["400"],
  indigo["500"],
  indigo["600"],
  indigo["700"],
  indigo["800"],
  indigo["900"],
  indigo["A400"],
  indigo["A700"],
  lightBlue["800"],
  lightBlue["900"],
  lightGreen["900"],
  lime["900"],
  pink["600"],
  pink["700"],
  pink["800"],
  pink["900"],
  pink["A700"],
  purple["400"],
  purple["500"],
  purple["600"],
  purple["700"],
  purple["800"],
  purple["900"],
  purple["A700"],
  red["700"],
  red["800"],
  red["900"],
  red["A700"],
  teal["700"],
  teal["800"],
  teal["900"],
];

const hashCode = (str: string) => {
  let h = 0;
  for (let i = 0; i < str.length; i++)
    h = (Math.imul(31, h) + str.charCodeAt(i)) | 0;

  return Math.abs(h);
};

const ID_HEX_COLOR_CACHE = new Map<string, string>();

/**
 * Convert ID to color
 * @param id ID
 * @returns color corresponding to id
 */
export const convertIdToHexColor = (id: string) => {
  if (ID_HEX_COLOR_CACHE.has(id)) return ID_HEX_COLOR_CACHE.get(id)!;
  const hexColor = COLOR_LIST[hashCode(id) % COLOR_LIST.length];
  ID_HEX_COLOR_CACHE.set(id, hexColor);
  return hexColor;
};
