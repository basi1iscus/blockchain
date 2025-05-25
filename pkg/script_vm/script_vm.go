package script_vm

import (
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/utils"
	"blockchain_demo/pkg/utils/queue"
	"blockchain_demo/pkg/utils/stack"
	"bufio"
	"encoding/hex"
	"errors"
	"strings"

	"fmt"
	"slices"
)

type OPCode byte

const (
	OP_FALSE = 0x00 // OP_0 is the opcode for 0
	OP_TRUE  = 0x01 // OP_1 is the opcode for 1
)
const (
	OP_0            = 0x00 // OP_0 is the opcode for 0
	OP_PUSHDATA     = 0x01 // OP_1 is the opcode for 1
	OP_PUSHDATA_4B  = 0x4B // OP_1 is the opcode for 1
	OP_PUSHDATA1    = 0x4C // OP_1 is the opcode for 1
	OP_PUSHDATA2    = 0x4D // OP_1 is the opcode for 1
	OP_PUSHDATA4    = 0x4E // OP_1 is the opcode for 1
	OP_1NEGATE      = 0x4F // OP_1NEGATE is the opcode for -1

	OP_1  = 0x51 // OP_1 is the opcode for 1
	OP_2  = 0x52 // OP_2 is the opcode for 2
	OP_3  = 0x53 // OP_3 is the opcode for 3
	OP_4  = 0x54 // OP_4 is the opcode for 4
	OP_5  = 0x55 // OP_5 is the opcode for 5
	OP_6  = 0x56 // OP_6 is the opcode for 6
	OP_7  = 0x57 // OP_7 is the opcode for 7
	OP_8  = 0x58 // OP_8 is the opcode for 8
	OP_9  = 0x59 // OP_9 is the opcode for 9
	OP_10 = 0x5A // OP_10 is the opcode for 10
	OP_11 = 0x5B // OP_11 is the opcode for 11
	OP_12 = 0x5C // OP_12 is the opcode for 12
	OP_13 = 0x5D // OP_13 is the opcode for 13
	OP_14 = 0x5E // OP_14 is the opcode for 14
	OP_15 = 0x5F // OP_15 is the opcode for 15
	OP_16 = 0x60 // OP_16 is the opcode for 16

	OP_NOP    = 0x61 //	Не делает ничего	Активен
	OP_IF     = 0x63 //	Выполняет следующие утверждения если верхнее значение стека не равно 0	Активен
	OP_NOTIF  = 0x64 //	Выполняет следующие утверждения если верхнее значение стека равно 0	Активен
	OP_ELSE   = 0x67 //	Выполняет утверждения если предыдущее OP_IF или OP_NOTIF было ложным	Активен
	OP_ENDIF  = 0x68 //	Завершает блок OP_IF, OP_NOTIF или OP_ELSE	Активен
	OP_VERIFY = 0x69 //	Проверяет верхнее значение стека и прерывает, если оно равно 0	Активен
	OP_RETURN = 0x6A //	Завершает выполнение и помечает транзакцию как недействительную

	OP_IFDUP        = 0x73 //	Дублирует верхний элемент стека, если он не равен 0	Активен
	OP_DEPTH        = 0x74 //	Помещает размер стека в стек	Активен
	OP_DROP         = 0x75 //	Удаляет верхний элемент стека	Активен
	OP_DUP          = 0x76 //	Дублирует верхний элемент стека	Активен
	OP_NIP          = 0x77 //	Удаляет второй элемент стека	Активен
	OP_OVER         = 0x78 //	Копирует второй элемент стека на верх	Активен
	OP_PICK         = 0x79 //	Копирует N-й элемент стека на верх (N из вершины стека)	Активен
	OP_ROLL         = 0x7A //	Перемещает N-й элемент стека на верх	Активен
	OP_ROT          = 0x7B //	Перемещает третий элемент стека на верх	Активен
	OP_SWAP         = 0x7C //	Меняет местами два верхних элемента стека	Активен
	OP_TUCK         = 0x7D //	Копирует верхний элемент стека под второй элемент	Активен
	OP_TOALTSTACK   = 0x6B //	Перемещает верхний элемент в альтернативный стек	Активен
	OP_FROMALTSTACK = 0x6C //	Перемещает верхний элемент из альтернативного стека	Активен
	OP_2DROP        = 0x6D //	Удаляет два верхних элемента стека	Активен
	OP_2DUP         = 0x6E //	Дублирует два верхних элемента стека	Активен
	OP_3DUP         = 0x6F //	Дублирует три верхних элемента стека	Активен
	OP_2OVER        = 0x70 //	Копирует два элемента с вершин стека на верх	Активен
	OP_2ROT         = 0x71 //	Пятый и шестой элементы перемещаются на верх стека	Активен
	OP_2SWAP        = 0x72 //

	OP_BOOLAND            = 0x9A //	Логическое И двух верхних элементов стека	Активен
	OP_BOOLOR             = 0x9B //	Логическое ИЛИ двух верхних элементов стека	Активен
	OP_NUMEQUAL           = 0x9C //	Возвращает 1, если два верхних элемента равны, иначе 0	Активен
	OP_NUMEQUALVERIFY     = 0x9D //	OP_NUMEQUAL + OP_VERIFY	Активен
	OP_NUMNOTEQUAL        = 0x9E //	Возвращает 1, если два верхних элемента не равны, иначе 0	Активен
	OP_LESSTHAN           = 0x9F //	Возвращает 1, если второй сверху < верхнего, иначе 0	Активен
	OP_GREATERTHAN        = 0xA0 //	Возвращает 1, если второй сверху > верхнего, иначе 0	Активен
	OP_LESSTHANOREQUAL    = 0xA1 //	Возвращает 1, если второй сверху ≤ верхнего, иначе 0	Активен
	OP_GREATERTHANOREQUAL = 0xA2 //	Возвращает 1, если второй сверху ≥ верхнего, иначе 0	Активен
	OP_MIN                = 0xA3 //	Возвращает меньшее из двух верхних элементов	Активен
	OP_MAX                = 0xA4 //	Возвращает большее из двух верхних элементов	Активен
	OP_EQUAL              = 0x87 //	Возвращает 1, если два верхних элемента равны, иначе 0	Активен
	OP_EQUALVERIFY        = 0x88 //	OP_EQUAL + OP_VERIFY	Акти

	OP_SHA256    = 0xA8 // 	Вычисляет SHA-256 хеш верхнего элемента	Активен
	OP_HASH160   = 0xA9 // 	Вычисляет RIPEMD-160(SHA-256()) верхнего элемента	Активен
	OP_HASH256   = 0xAA // 	Вычисляет SHA-256(SHA-256()) верхнего элемента	Активен
	OP_RIPEMD160 = 0xA6 // 	Вычисляет RIPEMD-160 хеш верхнего элемента	Активен
	OP_SHA1      = 0xA7 // 	Вычисляет SHA-1 хеш верхнего элемента

	OP_CHECKSIG            = 0xAC //	Проверяет подпись и публичный ключ	Активен
	OP_CHECKSIGVERIFY      = 0xAD //	OP_CHECKSIG + OP_VERIFY	Активен
	OP_CHECKMULTISIG       = 0xAE //	Проверяет несколько подписей и публичных ключей	Активен
	OP_CHECKMULTISIGVERIFY = 0xAF //	OP_CHECKMULTISIG + OP_VERIFY
)

var OpCodeNames = map[OPCode]string{
	OP_0:                   "OP_0",
	OP_1NEGATE:             "OP_1NEGATE",
	OP_1:                   "OP_1",
	OP_2:                   "OP_2",
	OP_3:                   "OP_3",
	OP_4:                   "OP_4",
	OP_5:                   "OP_5",
	OP_6:                   "OP_6",
	OP_7:                   "OP_7",
	OP_8:                   "OP_8",
	OP_9:                   "OP_9",
	OP_10:                  "OP_10",
	OP_11:                  "OP_11",
	OP_12:                  "OP_12",
	OP_13:                  "OP_13",
	OP_14:                  "OP_14",
	OP_15:                  "OP_15",
	OP_16:                  "OP_16",
	OP_PUSHDATA:            "OP_PUSHDATA",
	OP_PUSHDATA_4B:         "OP_PUSHDATA_4B",
	OP_PUSHDATA1:           "OP_PUSHDATA1",
	OP_PUSHDATA2:           "OP_PUSHDATA2",
	OP_PUSHDATA4:           "OP_PUSHDATA4",
	OP_NOP:                 "OP_NOP",
	OP_IFDUP:               "OP_IFDUP",
	OP_DROP:                "OP_DROP",
	OP_DUP:                 "OP_DUP",
	OP_EQUAL:               "OP_EQUAL",
	OP_EQUALVERIFY:         "OP_EQUALVERIFY",
	OP_VERIFY:              "OP_VERIFY",
	OP_SHA256:              "OP_SHA256",
	OP_HASH160:             "OP_HASH160",
	OP_HASH256:             "OP_HASH256",
	OP_CHECKSIG:            "OP_CHECKSIG",
	OP_CHECKSIGVERIFY:      "OP_CHECKSIGVERIFY",
	OP_IF:                  "OP_IF",
	OP_NOTIF:               "OP_NOTIF",
	OP_ELSE:                "OP_ELSE",
	OP_ENDIF:               "OP_ENDIF",
	OP_RETURN:              "OP_RETURN",
	OP_CHECKMULTISIG:       "OP_CHECKMULTISIG",
	OP_CHECKMULTISIGVERIFY: "OP_CHECKMULTISIGVERIFY",
}

var NamesOpCode = utils.ReverseMap(OpCodeNames)

type operation struct{
	scriptCode OPCode
	code OPCode
	data []byte
	childBranch branch
}

type branch struct{
	queue *queue.Queue[operation]
	parent *branch
}

type VM struct {
	stack  *stack.Stack[[]byte]
	mainBranch  branch
	signer sign.Signer
}

func IsActive(op OPCode) bool {
	for key := range OpCodeNames {
		if op == key {
			return true
		}
	}
	return false
}

func New(signer sign.Signer) *VM {
	return &VM{
		stack:  stack.New[[]byte](),
		signer: signer,
		mainBranch: branch{queue: queue.New[operation](), parent: nil},
	}
}

func compare(a, b []byte) int {
	if len(b) > len(a) {
		return 1
	} else if len(b) < len(a) {
		return -1
	}
	for i := 0; i < len(b); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

func (v *VM) equal() (bool, error) {
	a, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	b, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	return compare(a, b) == 0, nil
}

func (v *VM) topTrue() (bool, []byte, error) {
	top, err := v.stack.Pop()
	if err != nil {
		return false, nil, err
	}
	return compare(top, make([]byte, len(top))) != 0, top, nil
}

func (v *VM) checksig(data []byte) (bool, error) {
	pubKey, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	signature, err := v.stack.Pop()
	if err != nil {
		return false, err
	}

	return v.signer.Verify(data, signature, pubKey)
}

func (v *VM) checkmultisig(data []byte) (bool, error) {
	count, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	pubkeys := make([][]byte, int(count[0]))
	for i := 0; i < int(count[0]); i++ {
		pubKey, err := v.stack.Pop()
		if err != nil {
			return false, err
		}
		pubkeys[i] = pubKey
	}
	need, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	for i := 0; i < int(need[0]); i++ {
		found := false
		signature, err := v.stack.Pop()
		if err != nil {
			return false, err
		}
		for k, pubKey := range pubkeys {
			ok, err := v.signer.Verify(data, signature, pubKey)
			if err == nil && ok {
				// remove the public key from the list to avoid use double signature
				pubkeys = slices.Delete(pubkeys, k, k+1)
				found = true
				break
			}
		}
		if !found {	
			return false, nil		
		}
	}
	return true, nil
}

func (v *VM) ParseScript(script []byte) error {
	pointer := 0
	currentBranch := &v.mainBranch
	for pointer < len(script) {
		inc := 1
		dataLength := 0
		opCode := OPCode(script[pointer])
		switch {
		case opCode == OP_0:
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: OP_PUSHDATA, data: script[pointer : pointer+1]})
		case opCode == OP_1NEGATE:
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: OP_PUSHDATA, data: script[pointer : pointer+1]})
		case opCode >= OP_1 && opCode <= OP_16:
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: OP_PUSHDATA, data: []byte{byte(opCode) - OP_1 + 1}})
		case opCode >= OP_PUSHDATA && opCode <= OP_PUSHDATA_4B:
			dataLength = int(opCode)
			inc = dataLength + 1
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: OP_PUSHDATA, data: script[pointer+1 : pointer+1+dataLength]})
		case opCode == OP_PUSHDATA1:
			dataLength = int(script[pointer+1])
			inc = dataLength + 2
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: OP_PUSHDATA, data: script[pointer+2 : pointer+2+dataLength]})
		case opCode == OP_PUSHDATA2:
			dataLength = int(script[pointer+1]) | int(script[pointer+2])<<8
			inc = dataLength + 3
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: OP_PUSHDATA, data: script[pointer+3 : pointer+3+dataLength]})
		case opCode == OP_PUSHDATA4:
			dataLength = int(script[pointer+1]) | int(script[pointer+2])<<8 | int(script[pointer+3])<<16 | int(script[pointer+4])<<24
			inc = dataLength + 5
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: OP_PUSHDATA, data: script[pointer+5 : pointer+5+dataLength]})
		case opCode == OP_IF || opCode == OP_NOTIF:
			child := branch{queue: queue.New[operation](), parent: currentBranch}
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: opCode, data: nil, childBranch: child})
			currentBranch = &child
		case opCode == OP_ELSE:
			parent := currentBranch.parent
			if parent == nil {
				return errors.New("else without if")
			}
			child := branch{queue: queue.New[operation](), parent: parent}
			parent.queue.Enqueue(operation{scriptCode: opCode, code: opCode, data: nil, childBranch: child})
			currentBranch = &child
		case opCode == OP_ENDIF:
			parent := currentBranch.parent
			if parent == nil {
				return errors.New("endif without if")
			}
			currentBranch = parent
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: opCode, data: nil})
		case IsActive(opCode):
			currentBranch.queue.Enqueue(operation{scriptCode: opCode, code: opCode, data: nil})
		default:
			return fmt.Errorf("unknown opcode %#x", opCode)
		}

		pointer += inc
	}
	return nil
}

func Compile(op OPCode, data []byte) ([]byte, error) {
	if op == OP_PUSHDATA && len(data) == 0 {
		return nil, errors.New("OP_PUSHDATA cannot be used with empty data")
	} else if op == OP_PUSHDATA && len(data) == 1 && data[0] <= 16 {
		// If data is a single byte between 0 and 16, use the corresponding OP_1 to OP_16
		return []byte{byte(OP_1 + data[0] - 1)}, nil
	} else if op == OP_PUSHDATA && len(data) <= OP_PUSHDATA_4B {
		return append([]byte{byte(len(data))}, data...), nil
	} else if (op == OP_PUSHDATA || op == OP_PUSHDATA1) && len(data) <= 0xFF {
		return append([]byte{OP_PUSHDATA1, byte(len(data))}, data...), nil
	} else if (op == OP_PUSHDATA || op == OP_PUSHDATA2) && len(data) <= 0xFFFF {
		return append([]byte{OP_PUSHDATA2, byte(len(data)), byte(len(data) >> 8)}, data...), nil
	} else if (op == OP_PUSHDATA4 || op == OP_PUSHDATA2) && len(data) <= 0xFFFFFFFF {
		return append([]byte{OP_PUSHDATA4, byte(len(data)), byte(len(data) >> 8), byte(len(data) >> 16), byte(len(data) >> 24)}, data...), nil
	} else if op == OP_PUSHDATA || op == OP_PUSHDATA1 || op == OP_PUSHDATA2 || op == OP_PUSHDATA4 {
		return nil, fmt.Errorf("data length %d exceeds maximum allowed for opcode %s", len(data), OpCodeNames[op])
	}
	return []byte{byte(op)}, nil
}

func (v *VM) ParseString(s string) ([]byte, error) {
	row := 1
	script := make([]byte, 0)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}
		result := strings.Split(line, " ")
		op := result[0]
		var data []byte
		if len(result) > 1 {
			var err error
			data, err = hex.DecodeString(result[1])
			if err != nil {
				return nil, fmt.Errorf("failed to decode hex data in row %v: %w", row, err)
			}
		}
		opCode, ok := NamesOpCode[op]
		if !ok {	
			return nil, fmt.Errorf("failed to decode hex data in row %v: %s", row, op)
		}
		scriptOp, err := Compile(opCode, data)
		if err != nil {
			return nil, fmt.Errorf("failed to compile opcode %s in row %v: %w", op, row, err)
		}
		script = append(script, scriptOp...)
		row++
	}	
	return script, nil
}

func (v *VM) Run(script []byte, signedData []byte) ([]byte, error) {
	err := v.ParseScript(script)
	if err != nil {
		return nil, err
	}
	res, err := v.Execute(signedData, nil)
	if err != nil {	
		return res, err
	}
	return res, nil
}

func (v *VM) Execute(signedData []byte, workBranch *branch) ([]byte, error) {
	branch := workBranch
	if workBranch == nil {
		branch = &v.mainBranch
	}
	var lastIf bool
	for op := range branch.queue.Iterator() {
		switch (op.code) {
		case OP_PUSHDATA:
			v.stack.Push(op.data)
		case OP_IF:
			var err error
			lastIf, _, err = v.topTrue()
			if err == nil && lastIf {
				_, err = v.Execute(signedData, &op.childBranch)
				if err != nil {
					return nil, err
				}
			} else if err != nil {
				return nil, err
			}
		case OP_NOTIF:
			var err error
			lastIf, _, err = v.topTrue()
			lastIf = !lastIf
			if err == nil && lastIf {
				_, err = v.Execute(signedData, &op.childBranch)
				if err != nil {
					return nil, err
				}
			} else if err != nil {
				return nil, err
			}
		case OP_ELSE:
			if !lastIf {
				_, err := v.Execute(signedData, &op.childBranch)
				if err != nil {
					return nil, err
				}
			}
		case OP_ENDIF:
		case OP_NOP:
		case OP_RETURN:
			return nil, errors.New("return opcode encountered")
		case OP_IFDUP:
			top, err := v.stack.Pick()
			if err != nil {
				return nil, err
			}
			if compare(top, make([]byte, len(top))) != 0 {
				v.stack.Push(top)
			}
		case OP_DROP:
			_, err := v.stack.Pop()
			if err != nil {
				return nil, err
			}
		case OP_DUP:
			top, err := v.stack.Pick()
			if err != nil {
				return nil, err
			}
			v.stack.Push(top)
		case OP_EQUAL:
			eq, err := v.equal()
			if err != nil {
				return nil, err
			}
			if eq {
				v.stack.Push([]byte{OP_TRUE})
			} else {
				v.stack.Push([]byte{OP_FALSE})
			}
		case OP_EQUALVERIFY:
			eq, err := v.equal()
			if err != nil {
				return nil, err
			}
			if !eq {
				return nil, errors.New("equal verify failed")
			}
		case OP_VERIFY:
			if ok, _, err := v.topTrue(); err == nil || !ok {
				return nil, fmt.Errorf("verify failed")
			}			
		case OP_SHA256:
			top, err := v.stack.Pop()
			if err != nil {
				return nil, err
			}
			hash, err := utils.GetHash(top)
			if err != nil {
				return nil, err
			}
			v.stack.Push(hash)
		case OP_HASH160:
			top, err := v.stack.Pop()
			if err != nil {
				return nil, err
			}
			hash, err := utils.GetHash(top)
			if err != nil {
				return nil, err
			}
			hash160, err := utils.GetHash160(nil, hash)
			if err != nil {
				return nil, err
			}
			v.stack.Push(hash160)
		case OP_HASH256:
			top, err := v.stack.Pop()
			if err != nil {
				return nil, err
			}
			hash, err := utils.GetHash(top)
			if err != nil {
				return nil, err
			}
			hash256, err := utils.GetHash(hash)
			if err != nil {
				return nil, err
			}
			v.stack.Push(hash256)
		case OP_CHECKSIG:
			ok, err := v.checksig(signedData)
			if err != nil {
				return nil, err
			}
			if ok {
				v.stack.Push([]byte{OP_TRUE})
			} else {
				v.stack.Push([]byte{OP_FALSE})
			}
		case OP_CHECKSIGVERIFY:
			ok, err := v.checksig(signedData)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, errors.New("checksig verify failed")
			}
		case OP_CHECKMULTISIG:
			ok, err := v.checkmultisig(signedData)
			if err != nil {
				return nil, err
			}
			if ok {
				v.stack.Push([]byte{OP_TRUE})
			} else {
				v.stack.Push([]byte{OP_FALSE})
			}
		case OP_CHECKMULTISIGVERIFY:
			ok, err := v.checkmultisig(signedData)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, fmt.Errorf("checksig verify failed")
			}
		}
	}

	if workBranch != nil {
		return nil, nil
	}
	ok, top, err := v.topTrue()

	if err != nil || !ok {
		return top, fmt.Errorf("top of stack is not true, execution failed")
	}
	
	return top, nil
}

func (v *VM) String() string {

	branch := v.mainBranch.queue.ToArray()
	var result string
	for _, op := range branch {
		name, ok := OpCodeNames[op.code]
		if !ok {
			name = "UNKNOWN"
		}
		result += fmt.Sprintf("%s %x #0x%X\n", name, op.data, op.scriptCode)
	}
	
	return result
}

func (v *VM) GetStack() []string {
	result := make([]string, 0, v.stack.Size())
	stack := v.stack.ToArray()
	for _, op := range stack {
		result = append(result, hex.EncodeToString(op))
	}
	
	return result
}
