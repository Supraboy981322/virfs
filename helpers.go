package virfs

import ("strings")

func Is_valid_path(path string) bool {
	if len(path) < 1 { return false }
	if path[0] != '/' { return false }
	return true
}

func Get_basepath(path string) (string, error) {
	if !Is_valid_path(path) { return "", InvalidPath }
	if len(path) < 2 { return "/", nil }
	split := strings.Split(path[1:], "/")
	return strings.Join(split[:len(split)-1], "/"), nil
}

func (fs Fs) goto_path(path string) (*Dir, error) {
	if !Is_valid_path(path) { return nil, InvalidPath }

	base, e := Get_basepath(path)
	if len(base) == 0 { return &fs.Root, nil }
	if e != nil { return nil, e }
	path_split := strings.Split(base, "/")

	current := &(fs.Root)
	for _, d := range path_split {
		if !current.Contains(d) { return nil, DirNotExist }
		current = (*current).Content[d].Dir
	}
	return current, nil
}

func Get_name(path string) string {
	split := strings.Split(path, "/") 
	return split[len(split)-1]
}

func (d Dir) Contains(name string) bool {
	for n := range d.Content {
		if name == n { return true }
	}
	return false
}
